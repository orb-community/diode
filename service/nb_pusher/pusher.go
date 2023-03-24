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
	"github.com/netbox-community/go-netbox/v3/netbox/models"
	"github.com/orb-community/diode/service/config"
	"go.uber.org/zap"
)

type Pusher interface {
	Start() error
	Stop() error
}

type NetboxPusher struct {
	ctx    context.Context
	logger *zap.Logger
	config *config.Config
	client *client.NetBoxAPI
}

var _ Pusher = (*NetboxPusher)(nil)

var (
	unknown_name   string = "Unknown"
	unknown_slug   string = "unknown"
	staging_status string = "staging"
)

func New(ctx context.Context, logger *zap.Logger, config *config.Config) Pusher {
	return &NetboxPusher{ctx: ctx, logger: logger, config: config}
}

func (nb *NetboxPusher) Start() error {
	t := transport.New(nb.config.NetboxPusher.Endpoint, client.DefaultBasePath, []string{"https", "http"})
	t.DefaultAuthentication = transport.APIKeyAuth(
		"Authorization",
		"header",
		fmt.Sprintf("Token %v", nb.config.NetboxPusher.Token),
	)
	nb.client = client.New(t, nil)
	unkSiteCheck := dcim.NewDcimSitesListParams()
	unkSiteCheck.Slug = &unknown_slug
	var err error
	var list *dcim.DcimSitesListOK
	list, err = nb.client.Dcim.DcimSitesList(unkSiteCheck, nil)
	if err != nil {
		return err
	}
	if *list.GetPayload().Count != 0 {
		nb.logger.Info("netbox pusher started")
		return nil
	}

	nb.logger.Info("default unkown site does not exist, creating it")

	unkSite := dcim.NewDcimSitesCreateParams()
	unkSite.Data = &models.WritableSite{
		Name:   &unknown_name,
		Slug:   &unknown_slug,
		Status: staging_status,
	}

	_, err = nb.client.Dcim.DcimSitesCreate(unkSite, nil)
	if err != nil {
		nb.logger.Error("error", zap.Error(err))
		return err
	}
	nb.logger.Info("unknown site created")
	nb.logger.Info("netbox pusher started")
	return nil
}

func (nb *NetboxPusher) Stop() error {
	return nil
}
