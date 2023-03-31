/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package nb_pusher

import (
	"context"
	"encoding/json"
	"fmt"

	transport "github.com/go-openapi/runtime/client"
	"github.com/gosimple/slug"
	"github.com/netbox-community/go-netbox/v3/netbox/client"
	"github.com/netbox-community/go-netbox/v3/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/v3/netbox/client/extras"
	"github.com/netbox-community/go-netbox/v3/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/v3/netbox/client/status"
	"github.com/netbox-community/go-netbox/v3/netbox/models"
	"github.com/orb-community/diode/service/config"

	"go.uber.org/zap"
)

type Pusher interface {
	Start() error
	Stop() error
	CreateDevice([]byte) (int64, error)
	CreateInterface([]byte) (int64, error)
	CreateInterfaceIpAddress([]byte) (int64, error)
}

type NetboxPusher struct {
	ctx            context.Context
	logger         *zap.Logger
	config         *config.Config
	client         *client.NetBoxAPI
	unkSiteID      int64
	unkRoleID      int64
	unkDtypeID     int64
	unkMfrID       int64
	unkPlatID      int64
	tagsInit       bool
	discoveryTag   []*models.NestedTag
	placeholderTag []*models.NestedTag
}

var _ Pusher = (*NetboxPusher)(nil)

var (
	unknown_name           string = "Unknown"
	unknown_slug           string = slug.Make(unknown_name)
	discovery_tag_name     string = "Discovered"
	discovery_tag_slug     string = slug.Make(discovery_tag_name)
	placeholder_tag_name   string = "Placeholder"
	placeholder_tag_slug   string = slug.Make(placeholder_tag_name)
	unknown_interface_type string = "other"
)

const (
	invalid_id            int64  = -1
	staging_status        string = "staging"
	discovery_tag_color   string = "c0c0c0"
	placeholder_tag_color string = "ff6600"
)

func New(ctx context.Context, logger *zap.Logger, config *config.Config) Pusher {
	return &NetboxPusher{ctx: ctx, logger: logger, config: config, unkSiteID: invalid_id,
		unkRoleID: invalid_id, unkDtypeID: invalid_id, unkMfrID: invalid_id, tagsInit: false}
}

func (nb *NetboxPusher) Start() error {
	t := transport.New(nb.config.NetboxPusher.Endpoint, client.DefaultBasePath, []string{nb.config.NetboxPusher.Protocol})
	t.DefaultAuthentication = transport.APIKeyAuth(
		"Authorization",
		"header",
		fmt.Sprintf("Token %v", nb.config.NetboxPusher.Token),
	)
	nb.client = client.New(t, nil)
	if _, err := nb.client.Status.StatusList(status.NewStatusListParams(), nil); err != nil {
		return err
	}
	nb.logger.Info("netbox pusher started")
	return nil
}

func (nb *NetboxPusher) Stop() error {
	return nil
}

func (nb *NetboxPusher) CreateDevice(j []byte) (int64, error) {
	var err error
	if !nb.tagsInit {
		if err = nb.initializeDiodeTags(); err != nil {
			return invalid_id, err
		}
	}
	var deviceData NetboxDevice
	if err = json.Unmarshal(j, &deviceData); err != nil {
		return invalid_id, err
	}

	device := dcim.NewDcimDevicesCreateParams()
	var data models.WritableDeviceWithConfigContext

	var siteID int64
	if deviceData.Site != nil {
		deviceData.Site.Slug = slug.Make(deviceData.Site.Name)
		siteID, err = nb.createSite(deviceData.Site, nb.discoveryTag)
		if err != nil {
			return invalid_id, err
		}
		data.Site = &siteID
	} else {
		if nb.unkSiteID == invalid_id {
			if nb.unkSiteID, err = nb.createSite(&NetboxSite{Name: unknown_name, Slug: unknown_slug, Status: staging_status}, nb.placeholderTag); err != nil {
				return invalid_id, err
			}
		}
		data.Site = &nb.unkSiteID
	}

	var roleID int64
	if deviceData.Role != nil {
		deviceData.Role.Slug = slug.Make(deviceData.Role.Name)
		roleID, err = nb.createDeviceRole(deviceData.Role, nb.discoveryTag)
		if err != nil {
			return invalid_id, err
		}
		data.DeviceRole = &roleID
	} else {
		if nb.unkRoleID == invalid_id {
			unkownObject := &NetboxObject{Name: unknown_name, Slug: unknown_slug}
			if nb.unkRoleID, err = nb.createDeviceRole(unkownObject, nb.placeholderTag); err != nil {
				return invalid_id, err
			}
		}
		data.DeviceRole = &nb.unkRoleID
	}

	var typeID int64
	if deviceData.Type != nil {
		deviceData.Type.Slug = slug.Make(deviceData.Type.Model)
		typeID, err = nb.createDeviceType(deviceData.Type, nb.discoveryTag)
		if err != nil {
			return invalid_id, err
		}
		data.DeviceType = &typeID
	} else {
		if nb.unkDtypeID == invalid_id {
			if nb.unkDtypeID, err = nb.createDeviceType(&NetboxDeviceType{Mfr: nil, Model: unknown_name, Slug: unknown_slug}, nb.placeholderTag); err != nil {
				return invalid_id, err
			}
		}
		data.DeviceType = &nb.unkDtypeID
	}

	var platID int64
	if deviceData.Platform != nil {
		deviceData.Platform.Slug = slug.Make(deviceData.Platform.Name)
		platID, err = nb.createPlatform(deviceData.Platform, nb.discoveryTag)
		if err != nil {
			return invalid_id, err
		}
		data.Platform = &platID
	} else {
		if nb.unkDtypeID == invalid_id {
			if nb.unkPlatID, err = nb.createPlatform(&NetboxPlatform{Mfr: nil, Name: unknown_name, Slug: unknown_slug}, nb.placeholderTag); err != nil {
				return invalid_id, err
			}
		}
		data.Platform = &nb.unkPlatID
	}

	data.Status = DeviceStatusMap[deviceData.Status]
	data.Name = &deviceData.Name
	data.Serial = deviceData.Serial
	data.Tags = nb.discoveryTag

	device.Data = &data
	var created *dcim.DcimDevicesCreateCreated
	created, err = nb.client.Dcim.DcimDevicesCreate(device, nil)
	if err != nil {
		return invalid_id, err
	}
	nb.logger.Info("device created", zap.String("device", deviceData.Name))
	return created.Payload.ID, nil
}

func (nb *NetboxPusher) CreateInterface(j []byte) (int64, error) {
	var err error
	if !nb.tagsInit {
		if err = nb.initializeDiodeTags(); err != nil {
			return invalid_id, err
		}
	}
	var interfaceData NetboxInterface
	if err = json.Unmarshal(j, &interfaceData); err != nil {
		return invalid_id, err
	}

	ifs := dcim.NewDcimInterfacesCreateParams()
	var data models.WritableInterface

	data.Device = &interfaceData.DeviceID
	data.Name = &interfaceData.Name
	data.Vdcs = []int64{}
	data.TaggedVlans = []int64{}
	data.WirelessLans = []int64{}
	if interfaceData.Mtu > INTERFACE_MTU_MIN {
		data.Mtu = &interfaceData.Mtu
	}
	if interfaceData.Speed < INTERFACE_SPEED_MAX {
		data.Speed = &interfaceData.Speed
	}
	data.MacAddress = &interfaceData.MacAddress
	data.Enabled = InterfaceStateMap[interfaceData.State]
	data.Type = &unknown_interface_type
	data.Description = interfaceData.Type
	data.Tags = nb.discoveryTag

	ifs.Data = &data
	var created *dcim.DcimInterfacesCreateCreated
	created, err = nb.client.Dcim.DcimInterfacesCreate(ifs, nil)
	if err != nil {
		return invalid_id, err
	}
	nb.logger.Info("interface created", zap.String("interface", interfaceData.Name))
	return created.Payload.ID, nil
}

func (nb *NetboxPusher) CreateInterfaceIpAddress(j []byte) (int64, error) {
	var err error
	if !nb.tagsInit {
		if err = nb.initializeDiodeTags(); err != nil {
			return invalid_id, err
		}
	}
	var ipData NetboxIpAddress
	if err = json.Unmarshal(j, &ipData); err != nil {
		return invalid_id, err
	}

	ip := ipam.NewIpamIPAddressesCreateParams()
	var data models.WritableIPAddress

	data.Address = &ipData.Address
	data.AssignedObjectID = &ipData.AsgdObjID
	data.AssignedObjectType = &INTERFACE_OBJ_TYPE
	data.Tags = nb.discoveryTag

	ip.Data = &data
	var created *ipam.IpamIPAddressesCreateCreated
	created, err = nb.client.Ipam.IpamIPAddressesCreate(ip, nil)
	if err != nil {
		return invalid_id, err
	}
	nb.logger.Info("ip address for interface created", zap.String("ip_address", ipData.Address))
	return created.Payload.ID, nil
}

func (nb *NetboxPusher) initializeDiodeTags() error {
	var err error
	if nb.discoveryTag, err = nb.createDiodeTag(&discovery_tag_name, &discovery_tag_slug, discovery_tag_color); err != nil {
		return err
	}
	if nb.placeholderTag, err = nb.createDiodeTag(&placeholder_tag_name, &placeholder_tag_slug, placeholder_tag_color); err != nil {
		return err
	}
	nb.placeholderTag = append(nb.placeholderTag, nb.discoveryTag...)
	nb.tagsInit = true
	return nil
}

func (nb *NetboxPusher) createDiodeTag(name *string, slug *string, color string) ([]*models.NestedTag, error) {
	tagCheck := extras.NewExtrasTagsListParams()
	tagCheck.Slug = slug

	var err error
	var list *extras.ExtrasTagsListOK
	list, err = nb.client.Extras.ExtrasTagsList(tagCheck, nil)
	if err != nil {
		return nil, err
	}
	if *list.GetPayload().Count == 0 {
		extraTag := extras.NewExtrasTagsCreateParams()
		extraTag.Data = &models.Tag{
			Name:  name,
			Slug:  slug,
			Color: color,
		}
		_, err = nb.client.Extras.ExtrasTagsCreate(extraTag, nil)
		if err != nil {
			return nil, err
		}
	}
	discoveryTag := []*models.NestedTag{
		{
			Name: name,
			Slug: slug,
		},
	}
	return discoveryTag, nil
}

func (nb *NetboxPusher) createSite(site *NetboxSite, tag []*models.NestedTag) (int64, error) {
	siteCheck := dcim.NewDcimSitesListParams()
	siteCheck.Slug = &site.Slug
	var err error
	var list *dcim.DcimSitesListOK
	list, err = nb.client.Dcim.DcimSitesList(siteCheck, nil)
	if err != nil {
		return invalid_id, err
	}
	if *list.GetPayload().Count != 0 {
		for _, result := range list.GetPayload().Results {
			//return first match
			return result.ID, nil
		}
	}
	newSite := dcim.NewDcimSitesCreateParams()

	newSite.Data = &models.WritableSite{
		Name:   &site.Name,
		Slug:   &site.Slug,
		Status: site.Status,
		Tags:   tag,
	}
	var created *dcim.DcimSitesCreateCreated
	created, err = nb.client.Dcim.DcimSitesCreate(newSite, nil)
	if err != nil {
		return invalid_id, err
	}
	nb.logger.Info("site created", zap.String("site", site.Name))
	return created.Payload.ID, nil
}

func (nb *NetboxPusher) createDeviceRole(role *NetboxObject, tag []*models.NestedTag) (int64, error) {
	roleCheck := dcim.NewDcimDeviceRolesListParams()
	roleCheck.Slug = &role.Slug
	var err error
	var list *dcim.DcimDeviceRolesListOK
	list, err = nb.client.Dcim.DcimDeviceRolesList(roleCheck, nil)
	if err != nil {
		return invalid_id, err
	}
	if *list.GetPayload().Count != 0 {
		for _, result := range list.GetPayload().Results {
			//return first match
			return result.ID, nil
		}
	}
	newRole := dcim.NewDcimDeviceRolesCreateParams()
	newRole.Data = &models.DeviceRole{
		Name: &role.Name,
		Slug: &role.Slug,
		Tags: tag,
	}
	var created *dcim.DcimDeviceRolesCreateCreated
	created, err = nb.client.Dcim.DcimDeviceRolesCreate(newRole, nil)
	if err != nil {
		return invalid_id, err
	}
	nb.logger.Info("device role created", zap.String("role", role.Name))
	return created.Payload.ID, nil
}

func (nb *NetboxPusher) createManufacturer(mfr *NetboxObject, tag []*models.NestedTag) (int64, error) {
	mfrCheck := dcim.NewDcimManufacturersListParams()
	mfrCheck.Slug = &mfr.Slug
	var err error
	var list *dcim.DcimManufacturersListOK
	list, err = nb.client.Dcim.DcimManufacturersList(mfrCheck, nil)
	if err != nil {
		return invalid_id, err
	}
	if *list.GetPayload().Count != 0 {
		for _, result := range list.GetPayload().Results {
			//return first match
			return result.ID, nil
		}
	}
	newMfr := dcim.NewDcimManufacturersCreateParams()
	newMfr.Data = &models.Manufacturer{
		Name: &mfr.Name,
		Slug: &mfr.Slug,
		Tags: tag,
	}
	var created *dcim.DcimManufacturersCreateCreated
	created, err = nb.client.Dcim.DcimManufacturersCreate(newMfr, nil)
	if err != nil {
		return invalid_id, err
	}
	nb.logger.Info("manufacturer type created", zap.String("manufacturer", mfr.Name))
	return created.Payload.ID, nil
}

func (nb *NetboxPusher) createDeviceType(dType *NetboxDeviceType, tag []*models.NestedTag) (int64, error) {
	dTypeCheck := dcim.NewDcimDeviceTypesListParams()
	dTypeCheck.Slug = &dType.Slug
	var err error
	var list *dcim.DcimDeviceTypesListOK
	list, err = nb.client.Dcim.DcimDeviceTypesList(dTypeCheck, nil)
	if err != nil {
		return invalid_id, err
	}
	if *list.GetPayload().Count != 0 {
		for _, result := range list.GetPayload().Results {
			//return first match
			return result.ID, nil
		}
	}

	var mfrID int64
	if dType.Mfr != nil {
		dType.Mfr.Slug = slug.Make(dType.Mfr.Name)
		mfrID, err = nb.createManufacturer(dType.Mfr, nb.discoveryTag)
		if err != nil {
			return invalid_id, err
		}
	} else {
		if nb.unkMfrID == invalid_id {
			unkownObject := &NetboxObject{Name: unknown_name, Slug: unknown_slug}
			if nb.unkMfrID, err = nb.createManufacturer(unkownObject, nb.placeholderTag); err != nil {
				return invalid_id, err
			}
		}
		mfrID = nb.unkMfrID
	}

	newDType := dcim.NewDcimDeviceTypesCreateParams()
	newDType.Data = &models.WritableDeviceType{
		Model:        &dType.Model,
		Slug:         &dType.Slug,
		Manufacturer: &mfrID,
		Tags:         tag,
	}
	var created *dcim.DcimDeviceTypesCreateCreated
	created, err = nb.client.Dcim.DcimDeviceTypesCreate(newDType, nil)
	if err != nil {
		return invalid_id, err
	}
	nb.logger.Info("device type created", zap.String("type", dType.Model))
	return created.Payload.ID, nil
}

func (nb *NetboxPusher) createPlatform(plat *NetboxPlatform, tag []*models.NestedTag) (int64, error) {
	platCheck := dcim.NewDcimPlatformsListParams()
	platCheck.Slug = &plat.Slug
	var err error
	var list *dcim.DcimPlatformsListOK
	list, err = nb.client.Dcim.DcimPlatformsList(platCheck, nil)
	if err != nil {
		return invalid_id, err
	}
	if *list.GetPayload().Count != 0 {
		for _, result := range list.GetPayload().Results {
			//return first match
			return result.ID, nil
		}
	}

	var mfrID int64
	if plat.Mfr != nil {
		plat.Mfr.Slug = slug.Make(plat.Mfr.Name)
		mfrID, err = nb.createManufacturer(plat.Mfr, nb.discoveryTag)
		if err != nil {
			return invalid_id, err
		}
	} else {
		if nb.unkMfrID == invalid_id {
			unkownObject := &NetboxObject{Name: unknown_name, Slug: unknown_slug}
			if nb.unkMfrID, err = nb.createManufacturer(unkownObject, nb.placeholderTag); err != nil {
				return invalid_id, err
			}
		}
		mfrID = nb.unkMfrID
	}

	newPlatform := dcim.NewDcimPlatformsCreateParams()
	newPlatform.Data = &models.WritablePlatform{
		Name:         &plat.Name,
		Slug:         &plat.Slug,
		Manufacturer: &mfrID,
		Tags:         tag,
	}
	var created *dcim.DcimPlatformsCreateCreated
	created, err = nb.client.Dcim.DcimPlatformsCreate(newPlatform, nil)
	if err != nil {
		return invalid_id, err
	}
	nb.logger.Info("device platform created", zap.String("platform", plat.Name))
	return created.Payload.ID, nil
}
