/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package nb_pusher

import (
	"context"
	"fmt"

	transport "github.com/go-openapi/runtime/client"
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
	Device([]byte) error
}

type NetboxPusher struct {
	ctx       context.Context
	logger    *zap.Logger
	config    *config.Config
	client    *client.NetBoxAPI
	unkSiteID int64
	unkRoleID int64
	diodeTag  []*models.NestedTag
}

var _ Pusher = (*NetboxPusher)(nil)

var (
	unknown_name   string = "Unknown"
	unknown_slug   string = "unknown"
	diode_tag_name string = "Diode"
	diode_tag_slug string = "diode"
)

const (
	invalid_id      int64  = -1
	staging_status  string = "staging"
	diode_tag_color string = "ff6600"
)

func New(ctx context.Context, logger *zap.Logger, config *config.Config) Pusher {
	return &NetboxPusher{ctx: ctx, logger: logger, config: config, unkSiteID: invalid_id, unkRoleID: invalid_id}
}

func (nb *NetboxPusher) Start() error {
	t := transport.New(nb.config.NetboxPusher.Endpoint, client.DefaultBasePath, []string{"https", "http"})
	t.DefaultAuthentication = transport.APIKeyAuth(
		"Authorization",
		"header",
		fmt.Sprintf("Token %v", nb.config.NetboxPusher.Token),
	)
	nb.client = client.New(t, nil)

	var err error
	if nb.diodeTag, err = nb.createDiodeTag(); err != nil {
		return err
	}
	if nb.unkSiteID, err = nb.createUnknownSite(nb.diodeTag); err != nil {
		return err
	}
	if nb.unkRoleID, err = nb.createUnknownDeviceRole(nb.diodeTag); err != nil {
		return err
	}
	nb.logger.Info("netbox pusher started")
	return nil
}

func (nb *NetboxPusher) Stop() error {
	return nil
}

func (nb *NetboxPusher) Device(json []byte) error {
	return nil
}

func (nb *NetboxPusher) createDiodeTag() ([]*models.NestedTag, error) {
	diodeTagCheck := extras.NewExtrasTagsListParams()
	diodeTagCheck.Slug = &diode_tag_slug
	var err error
	var list *extras.ExtrasTagsListOK
	list, err = nb.client.Extras.ExtrasTagsList(diodeTagCheck, nil)
	if err != nil {
		return nil, err
	}
	if *list.GetPayload().Count == 0 {
		nb.logger.Info("default diode tag does not exist, creating it")
		diodeExtraTag := extras.NewExtrasTagsCreateParams()
		diodeExtraTag.Data = &models.Tag{
			Name:  &diode_tag_name,
			Slug:  &diode_tag_slug,
			Color: diode_tag_color,
		}
		_, err = nb.client.Extras.ExtrasTagsCreate(diodeExtraTag, nil)
		if err != nil {
			return nil, err
		}
	}
	diodeTag := make([]*models.NestedTag, 1)
	diodeTag[0] = &models.NestedTag{
		Name: &diode_tag_name,
		Slug: &diode_tag_slug,
	}
	return diodeTag, nil
}

func (nb *NetboxPusher) createUnknownSite(tag []*models.NestedTag) (int64, error) {
	unkSiteCheck := dcim.NewDcimSitesListParams()
	unkSiteCheck.Slug = &unknown_slug
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
	nb.logger.Info("default unkown site does not exist, creating it")
	unkSite := dcim.NewDcimSitesCreateParams()

	unkSite.Data = &models.WritableSite{
		Name:   &unknown_name,
		Slug:   &unknown_slug,
		Status: staging_status,
		Tags:   tag,
	}
	var created *dcim.DcimSitesCreateCreated
	created, err = nb.client.Dcim.DcimSitesCreate(unkSite, nil)
	if err != nil {
		return invalid_id, err
	}
	nb.logger.Info("unknown site created")
	return created.Payload.ID, nil
}

func (nb *NetboxPusher) createUnknownDeviceRole(tag []*models.NestedTag) (int64, error) {
	unkRoleCheck := dcim.NewDcimDeviceRolesListParams()
	unkRoleCheck.Slug = &unknown_slug
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
	nb.logger.Info("default unkown device role does not exist, creating it")
	unkRole := dcim.NewDcimDeviceRolesCreateParams()
	unkRole.Data = &models.DeviceRole{
		Name: &unknown_name,
		Slug: &unknown_slug,
		Tags: tag,
	}
	var created *dcim.DcimDeviceRolesCreateCreated
	created, err = nb.client.Dcim.DcimDeviceRolesCreate(unkRole, nil)
	if err != nil {
		return invalid_id, err
	}
	nb.logger.Info("unknown device role created")
	return created.Payload.ID, nil
}
