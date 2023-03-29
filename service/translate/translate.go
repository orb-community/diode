/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package translate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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

type ifIpJsonReturn struct {
	IfID int64  `json:"assigned_object_id"`
	Ip   string `json:"address"`
}

func New(ctx context.Context, logger *zap.Logger, config *config.Config, db *storage.Service, pusher *nb_pusher.Pusher) Translator {
	return &SuzieQTranslate{ctx: ctx, logger: logger, config: config, db: db, pusher: pusher}
}

func (st *SuzieQTranslate) Translate(data interface{}) error {
	if devices, ok := data.([]storage.DbDevice); ok {
		var errs error
		for _, device := range devices {
			if len(device.Id) == 0 {
				continue
			}
			j, err := st.translateDevice(&device)
			if err != nil {
				if errs != nil {
					errs = fmt.Errorf("%v; %v", errs, err)
				} else {
					errs = err
				}
				continue
			}
==== BASE ====
			_, err = (*st.pusher).CreateDevice(j)
==== BASE ====
			if err != nil {
				if errs != nil {
					errs = fmt.Errorf("%v; %v", errs, err)
				} else {
					errs = err
				}
				continue
			}
			newDevice, err := (*st.db).UpdateDevice(device.Id, id)
			if err != nil {
				if errs != nil {
					errs = fmt.Errorf("%v; %v", errs, err)
				} else {
					errs = err
				}
				continue
			}
			if err := st.checkExistingInterfaces(&newDevice); err != nil {
				if errs != nil {
					errs = fmt.Errorf("%v; %v", errs, err)
				} else {
					errs = err
				}
				continue
			}
			if err := st.checkExistingInterfaces(&newDevice); err != nil {
				return err
			}
		}
		return errs
	} else if ifs, ok := data.([]storage.DbInterface); ok {
		var errs error
		for _, ifce := range ifs {
			if len(ifce.Id) == 0 {
				continue
			}
			device, err := (*st.db).GetDeviceByPolicyAndNamespaceAndHostname(ifce.Policy, ifce.Namespace, ifce.Hostname)
			if err != nil {
				if errs != nil {
					errs = fmt.Errorf("%v; %v", errs, err)
				} else {
					errs = err
				}
				continue
			} else if device.NetboxRefId == invalid_id {
				err = errors.New("invalid device id")
				if errs != nil {
					errs = fmt.Errorf("%v; %v", errs, err)
				} else {
					errs = err
				}
				continue
			}
			j, err := st.translateInterface(&ifce, device.NetboxRefId)
			if err != nil {
				if errs != nil {
					errs = fmt.Errorf("%v; %v", errs, err)
				} else {
					errs = err
				}
				continue
			}
			id, err := (*st.pusher).CreateInterface(j)
			if err != nil {
				if errs != nil {
					errs = fmt.Errorf("%v; %v", errs, err)
				} else {
					errs = err
				}
				continue
			}
			newInterface, err := (*st.db).UpdateInterface(ifce.Id, id)
			if err != nil {
				if errs != nil {
					errs = fmt.Errorf("%v; %v", errs, err)
				} else {
					errs = err
				}
				continue
			}
			for _, ip := range newInterface.IpAddresses {
				j, err := st.translateIpInterface(&ip, newInterface.NetboxRefId)
				if err != nil {
					if errs != nil {
						errs = fmt.Errorf("%v; %v", errs, err)
					} else {
						errs = err
					}
					continue
				}
				if _, err := (*st.pusher).CreateInterfaceIpAddress(j); err != nil {
					if errs != nil {
						errs = fmt.Errorf("%v; %v", errs, err)
					} else {
						errs = err
					}
				}
			}
		}
		return errs
	} else if vlans, ok := data.([]storage.DbVlan); ok {
		var errs error
		for _, vlan := range vlans {
			if len(vlan.Id) == 0 {
				continue
			}
			if _, err := st.translateVlan(&vlan); err != nil {
				if errs != nil {
					errs = fmt.Errorf("%v; %v", errs, err)
				} else {
					errs = err
				}
				continue
			}
		}
		return errs
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

func (st *SuzieQTranslate) translateIpInterface(ip *storage.IpAddress, ifID int64) ([]byte, error) {
	var ret ifIpJsonReturn
	ret.IfID = ifID
	ret.Ip = ip.Address
	return json.Marshal(ret)
}

func (st *SuzieQTranslate) translateVlan(vlan *storage.DbVlan) ([]byte, error) {
	st.logger.Warn("translate for vlan data not implemented yet")
	return nil, nil
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
