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

const invalid_id int64 = -1

type SuzieQTranslate struct {
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

type ifJsonReturn struct {
	DeviceID   int64  `json:"device_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Speed      int64  `json:"speed"`
	Mtu        int64  `json:"mtu"`
	MacAddress string `json:"mac_address"`
	State      string `json:"state"`
}

func New(ctx context.Context, logger *zap.Logger, config *config.Config, db *storage.Service, pusher *nb_pusher.Pusher) Translator {
	return &SuzieQTranslate{ctx: ctx, logger: logger, config: config, db: db, pusher: pusher}
}

func (st *SuzieQTranslate) Translate(data interface{}) error {
	if devices, ok := data.([]storage.DbDevice); ok {
		for _, device := range devices {
			if len(device.Id) == 0 {
				continue
			}
			j, err := st.translateDevice(&device)
			if err != nil {
				return err
			}
			id, err := (*st.pusher).CreateDevice(j)
			if err != nil {
				return err
			}
			newDevice, err := (*st.db).UpdateDevice(device.Id, id)
			if err != nil {
				return err
			}
			if err := st.checkExistingInterfaces(&newDevice); err != nil {
				return err
			}
		}
		return nil
	} else if ifs, ok := data.([]storage.DbInterface); ok {
		for _, ifce := range ifs {
			if len(ifce.Id) == 0 {
				continue
			}
			device, err := (*st.db).GetDeviceByPolicyAndNamespaceAndHostname(ifce.Policy, ifce.Namespace, ifce.Hostname)
			if err != nil {
				return err
			} else if device.NetboxRefId == invalid_id {
				return errors.New("invalid device id")
			}
			j, err := st.translateInterface(&ifce, device.NetboxRefId)
			if err != nil {
				return err
			}
			id, err := (*st.pusher).CreateInterface(j)
			if err != nil {
				return err
			}
			if _, err := (*st.db).UpdateInterface(ifce.Id, id); err != nil {
				return err
			}
		}
		return nil
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

func (st *SuzieQTranslate) translateDevice(device *storage.DbDevice) ([]byte, error) {
	var ret deviceJsonReturn
	ret.Name = device.Hostname
	ret.Status = device.State
	ret.Serial = device.SerialNumber
	ret.Dtype.Model = device.Model
	ret.Dtype.Mfr.Name = device.Vendor
	return json.Marshal(ret)
}

func (st *SuzieQTranslate) translateInterface(ifs *storage.DbInterface, deviceID int64) ([]byte, error) {
	var ret ifJsonReturn
	ret.DeviceID = deviceID
	ret.Name = ifs.Name
	ret.Type = ifs.IfType
	ret.Speed = ifs.Speed
	ret.Mtu = ifs.Mtu
	ret.MacAddress = ifs.MacAddress
	ret.State = ifs.AdminState
	return json.Marshal(ret)
}

func (st *SuzieQTranslate) translateVlan(vlan *storage.DbVlan) ([]byte, error) {
	return nil, errors.New("translate for vlan data not implemented yet")
}

func (st *SuzieQTranslate) checkExistingInterfaces(device *storage.DbDevice) error {
	ifs, err := (*st.db).GetInterfaceByPolicyAndNamespaceAndHostname(device.Policy, device.Namespace, device.Hostname)
	if err != nil {
		return nil
	}
	var vIfs []storage.DbInterface

	for _, v := range ifs {
		if v.NetboxRefId == invalid_id {
			vIfs = append(vIfs, v)
		}
	}
	if len(vIfs) > 0 {
		return st.Translate(vIfs)
	}
	return nil
}
