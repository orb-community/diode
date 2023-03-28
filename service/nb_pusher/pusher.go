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
	"github.com/netbox-community/go-netbox/v3/netbox/models"
	"github.com/orb-community/diode/service/config"

	"go.uber.org/zap"
)

type Pusher interface {
	Start() error
	Stop() error
	CreateDevice([]byte) (int64, error)
	CreateInterface([]byte) (int64, error)
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
	tagsInit       bool
	discoveryTag   models.NestedTag
	placeholderTag models.NestedTag
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
	t := transport.New(nb.config.NetboxPusher.Endpoint, client.DefaultBasePath, []string{"https", "http"})
	t.DefaultAuthentication = transport.APIKeyAuth(
		"Authorization",
		"header",
		fmt.Sprintf("Token %v", nb.config.NetboxPusher.Token),
	)
	nb.client = client.New(t, nil)

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
		siteID, err = nb.createSite(deviceData.Site, []*models.NestedTag{&nb.discoveryTag})
		if err != nil {
			return invalid_id, err
		}
		data.Site = &siteID
	} else {
		if nb.unkSiteID == invalid_id {
			if nb.unkSiteID, err = nb.createSite(&NetboxSite{Name: unknown_name, Slug: unknown_slug, Status: staging_status},
				[]*models.NestedTag{&nb.discoveryTag, &nb.placeholderTag}); err != nil {
				return invalid_id, err
			}
		}
		data.Site = &nb.unkSiteID
	}

	var roleID int64
	if deviceData.Role != nil {
		deviceData.Role.Slug = slug.Make(deviceData.Role.Name)
		roleID, err = nb.createDeviceRole(deviceData.Role, []*models.NestedTag{&nb.discoveryTag})
		if err != nil {
			return invalid_id, err
		}
		data.DeviceRole = &roleID
	} else {
		if nb.unkRoleID == invalid_id {
			unkownObject := &NetboxObject{Name: unknown_name, Slug: unknown_slug}
			if nb.unkRoleID, err = nb.createDeviceRole(unkownObject, []*models.NestedTag{&nb.discoveryTag, &nb.placeholderTag}); err != nil {
				return invalid_id, err
			}
		}
		data.DeviceRole = &nb.unkRoleID
	}

	var typeID int64
	if deviceData.Type != nil {
		deviceData.Type.Slug = slug.Make(deviceData.Type.Model)
		typeID, err = nb.createDeviceType(deviceData.Type, []*models.NestedTag{&nb.discoveryTag})
		if err != nil {
			return invalid_id, err
		}
		data.DeviceType = &typeID
	} else {
		if nb.unkDtypeID == invalid_id {
			if nb.unkDtypeID, err = nb.createDeviceType(&NetboxDeviceType{Mfr: nil, Model: unknown_name, Slug: unknown_slug},
				[]*models.NestedTag{&nb.discoveryTag, &nb.placeholderTag}); err != nil {
				return invalid_id, err
			}
		}
		data.DeviceType = &nb.unkDtypeID
	}

	data.Status = DeviceStatusMap[deviceData.Status]
	data.Name = &deviceData.Name
	data.Tags = []*models.NestedTag{&nb.discoveryTag}

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

	ifs.Data.Device = &interfaceData.DeviceID
	ifs.Data.Name = &interfaceData.Name
	ifs.Data.Speed = &interfaceData.Speed
	ifs.Data.Mtu = &interfaceData.Mtu
	ifs.Data.MacAddress = &interfaceData.MacAddress
	ifs.Data.Enabled = InterfaceStateMap[interfaceData.State]
	ifs.Data.Type = &unknown_interface_type
	ifs.Data.Description = interfaceData.Type

	var created *dcim.DcimInterfacesCreateCreated
	created, err = nb.client.Dcim.DcimInterfacesCreate(ifs, nil)
	if err != nil {
		return invalid_id, err
	}

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
	nb.tagsInit = true
	return nil
}

func (nb *NetboxPusher) createDiodeTag(name *string, slug *string, color string) (models.NestedTag, error) {
	tagCheck := extras.NewExtrasTagsListParams()
	tagCheck.Slug = slug
	var discoveryTag models.NestedTag
	var err error
	var list *extras.ExtrasTagsListOK
	list, err = nb.client.Extras.ExtrasTagsList(tagCheck, nil)
	if err != nil {
		return discoveryTag, err
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
			return discoveryTag, err
		}
	}
	discoveryTag = models.NestedTag{
		Name: &discovery_tag_name,
		Slug: &discovery_tag_slug,
	}
	return discoveryTag, nil
}

func (nb *NetboxPusher) createSite(site *NetboxSite, tag []*models.NestedTag) (int64, error) {
	unkSiteCheck := dcim.NewDcimSitesListParams()
	unkSiteCheck.Slug = &site.Slug
	var err error
	var list *dcim.DcimSitesListOK
	list, err = nb.client.Dcim.DcimSitesList(unkSiteCheck, nil)
	if err != nil {
		return invalid_id, err
	}
	if *list.GetPayload().Count != 0 {
		for _, result := range list.GetPayload().Results {
			//return first match
			return result.ID, nil
		}
	}
	unkSite := dcim.NewDcimSitesCreateParams()

	unkSite.Data = &models.WritableSite{
		Name:   &site.Name,
		Slug:   &site.Slug,
		Status: site.Status,
		Tags:   tag,
	}
	var created *dcim.DcimSitesCreateCreated
	created, err = nb.client.Dcim.DcimSitesCreate(unkSite, nil)
	if err != nil {
		return invalid_id, err
	}
	nb.logger.Info("site created", zap.String("site", site.Name))
	return created.Payload.ID, nil
}

func (nb *NetboxPusher) createDeviceRole(role *NetboxObject, tag []*models.NestedTag) (int64, error) {
	unkRoleCheck := dcim.NewDcimDeviceRolesListParams()
	unkRoleCheck.Slug = &role.Name
	var err error
	var list *dcim.DcimDeviceRolesListOK
	list, err = nb.client.Dcim.DcimDeviceRolesList(unkRoleCheck, nil)
	if err != nil {
		return invalid_id, err
	}
	if *list.GetPayload().Count != 0 {
		for _, result := range list.GetPayload().Results {
			//return first match
			return result.ID, nil
		}
	}
	unkRole := dcim.NewDcimDeviceRolesCreateParams()
	unkRole.Data = &models.DeviceRole{
		Name: &role.Name,
		Slug: &role.Slug,
		Tags: tag,
	}
	var created *dcim.DcimDeviceRolesCreateCreated
	created, err = nb.client.Dcim.DcimDeviceRolesCreate(unkRole, nil)
	if err != nil {
		return invalid_id, err
	}
	nb.logger.Info("device role created", zap.String("role", role.Name))
	return created.Payload.ID, nil
}

func (nb *NetboxPusher) createManufacturer(mfr *NetboxObject, tag []*models.NestedTag) (int64, error) {
	unkMfrCheck := dcim.NewDcimManufacturersListParams()
	unkMfrCheck.Slug = &mfr.Name
	var err error
	var list *dcim.DcimManufacturersListOK
	list, err = nb.client.Dcim.DcimManufacturersList(unkMfrCheck, nil)
	if err != nil {
		return invalid_id, err
	}
	if *list.GetPayload().Count != 0 {
		for _, result := range list.GetPayload().Results {
			//return first match
			return result.ID, nil
		}
	}
	unkMfr := dcim.NewDcimManufacturersCreateParams()
	unkMfr.Data = &models.Manufacturer{
		Name: &mfr.Name,
		Slug: &mfr.Slug,
		Tags: tag,
	}
	var created *dcim.DcimManufacturersCreateCreated
	created, err = nb.client.Dcim.DcimManufacturersCreate(unkMfr, nil)
	if err != nil {
		return invalid_id, err
	}
	nb.logger.Info("manufacturer type created", zap.String("manufacturer", mfr.Name))
	return created.Payload.ID, nil
}

func (nb *NetboxPusher) createDeviceType(dType *NetboxDeviceType, tag []*models.NestedTag) (int64, error) {
	unkDTypeCheck := dcim.NewDcimDeviceTypesListParams()
	unkDTypeCheck.Slug = &dType.Slug
	var err error
	var list *dcim.DcimDeviceTypesListOK
	list, err = nb.client.Dcim.DcimDeviceTypesList(unkDTypeCheck, nil)
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
		mfrID, err = nb.createManufacturer(dType.Mfr, []*models.NestedTag{&nb.discoveryTag})
		if err != nil {
			return invalid_id, err
		}
	} else {
		if nb.unkMfrID == invalid_id {
			unkownObject := &NetboxObject{Name: unknown_name, Slug: unknown_slug}
			if nb.unkMfrID, err = nb.createManufacturer(unkownObject, []*models.NestedTag{&nb.discoveryTag, &nb.placeholderTag}); err != nil {
				return invalid_id, err
			}
		}
		mfrID = nb.unkMfrID
	}

	unkDType := dcim.NewDcimDeviceTypesCreateParams()
	unkDType.Data = &models.WritableDeviceType{
		Model:        &dType.Model,
		Slug:         &dType.Slug,
		Manufacturer: &mfrID,
		Tags:         tag,
	}
	var created *dcim.DcimDeviceTypesCreateCreated
	created, err = nb.client.Dcim.DcimDeviceTypesCreate(unkDType, nil)
	if err != nil {
		return invalid_id, err
	}
	nb.logger.Info("device type created", zap.String("type", dType.Model))
	return created.Payload.ID, nil
}
