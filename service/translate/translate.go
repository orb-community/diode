/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package translate

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/orb-community/diode/service/config"
	"github.com/orb-community/diode/service/storage"
	"go.uber.org/zap"
)

type Translator interface {
	Translate(interface{}) ([]byte, error)
}

type SuzieqTraslate struct {
	ctx    context.Context
	logger *zap.Logger
	config *config.Config
	db     *storage.Service
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

func New(ctx context.Context, logger *zap.Logger, config *config.Config, db *storage.Service) Translator {
	return &SuzieqTraslate{ctx: ctx, logger: logger, config: config, db: db}
}

func (st *SuzieqTraslate) Translate(data interface{}) ([]byte, error) {
	if device, ok := data.(storage.DbDevice); ok {
		return st.translateDevice(&device)
	} else if ifs, ok := data.(storage.DbInterface); ok {
		return st.translateInterface(&ifs)
	} else if vlan, ok := data.(storage.DbVlan); ok {
		return st.translateVlan(&vlan)
	}
	return nil, errors.New("no valid translatable data found")
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
