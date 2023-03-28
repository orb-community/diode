/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package translate

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/orb-community/diode/service/config"
	"github.com/orb-community/diode/service/nb_pusher"
	"github.com/orb-community/diode/service/storage"
	"go.uber.org/zap"
)

type Translator interface {
	Translate(interface{}) error
}

type SuzieqTraslate struct {
	ctx    context.Context
	logger *zap.Logger
	config *config.Config
	db     *storage.Service
	pusher *nb_pusher.Pusher
}

type deviceJsonReturn struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Serial string `json:"serial"`
	Dtype  struct {
		Model string `json:"model"`
		Mfr   struct {
			Name string `json:"name"`
		} `json:"manufacturers"`
	} `json:"device_type"`
}

func New(ctx context.Context, logger *zap.Logger, config *config.Config, db *storage.Service, pusher *nb_pusher.Pusher) Translator {
	return &SuzieqTraslate{ctx: ctx, logger: logger, config: config, db: db, pusher: pusher}
}

func (st *SuzieqTraslate) Translate(data interface{}) error {
	if devices, ok := data.([]storage.DbDevice); ok {
		for _, device := range devices {
			if len(device.Id) == 0 {
				continue
			}
			j, err := st.translateDevice(&device)
			if err != nil {
				return err
			}
			_, err = (*st.pusher).CreateDevice(j)
			if err != nil {
				return err
			}
		}
		return nil
	} else if ifs, ok := data.([]storage.DbInterface); ok {
		for _, ifce := range ifs {
			if len(ifce.Id) == 0 {
				continue
			}
			_, err := st.translateInterface(&ifce)
			return err
		}
	} else if vlans, ok := data.([]storage.DbVlan); ok {
		for _, vlan := range vlans {
			if len(vlan.Id) == 0 {
				continue
			}
			_, err := st.translateVlan(&vlan)
			return err
		}
	}
	return errors.New("no valid translatable data found")
}

func (st *SuzieqTraslate) translateDevice(device *storage.DbDevice) ([]byte, error) {
	var ret deviceJsonReturn
	ret.Name = device.Hostname
	ret.Status = device.State
	ret.Serial = device.SerialNumber
	ret.Dtype.Model = device.Model
	ret.Dtype.Mfr.Name = device.Vendor
	return json.Marshal(ret)
}

func (st *SuzieqTraslate) translateInterface(ifs *storage.DbInterface) ([]byte, error) {
	return nil, errors.New("translate for interface data not implemented yet")
}

func (st *SuzieqTraslate) translateVlan(vlan *storage.DbVlan) ([]byte, error) {
	return nil, errors.New("translate for vlan data not implemented yet")
}
