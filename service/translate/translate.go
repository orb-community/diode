/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package translate

import (
	"context"
	"encoding/json"
	"errors"
	"net"

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
	db     storage.Service
	pusher nb_pusher.Pusher
}

var ipChecker nb_pusher.NetboxPrimaryIpChecker

func New(ctx context.Context, logger *zap.Logger, config *config.Config, db storage.Service, pusher nb_pusher.Pusher) Translator {
	return &SuzieQTranslate{ctx: ctx, logger: logger, config: config, db: db, pusher: pusher}
}

func (st *SuzieQTranslate) Translate(data interface{}) error {
	if devices, ok := data.([]storage.DbDevice); ok {
		var errs error
		for _, device := range devices {

			if len(device.Id) == 0 {
				continue
			}

			// Separates host from port
			host, _, err := net.SplitHostPort(device.Address)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}
			
			var deviceAddresses []string

			device.Address = "12.12.12.1/24" // Change the device addr to be matched with one of the ifc
			deviceAddresses = append(deviceAddresses, device.Address)

			ipAddr := net.ParseIP(host) 
			if ipAddr != nil { 
				// Only enters here if the IP is valid
				device.Address = ipAddr.String()
			} else {
				_, err = net.LookupHost(host) 
				if err != nil {
					errs = errors.Join(errs, err)
					continue
				} else {
					ips, err := net.LookupIP(host)
					if err != nil {
						errs = errors.Join(errs, err)
						continue
					} else {
						for _, ip := range ips {
							deviceAddresses = append(deviceAddresses, ip.String())
							// A host may have multiple IPs
						}
						// Get the interfaces for that device to compare the IPs with the device IPs
						ifs, err := st.db.GetInterfaceByPolicyAndNamespaceAndHostname(device.Policy, device.Namespace, device.Hostname)
						if err != nil {
							errs = errors.Join(errs, err)
							continue
						}
						// If there are multiple IPs, then check if the address is in the list				
						for _, ifc := range ifs {
							if len(ifc.IpAddresses) > 0 {
								for idx := range ifc.IpAddresses {
									for k := range deviceAddresses {
										if ifc.IpAddresses[idx].Address == deviceAddresses[k] {
											st.logger.Info("matching ip addr between interface and device", zap.String("primary IP: ", deviceAddresses[k]))
											device.Address = ifc.IpAddresses[idx].Address 
											// Store the matched Ifc IP to the device.Address field (later will be translated and stored on the checker struct)
										}
									}
								}
							}
						}
					}
				}
			}

			j, err := st.translateDevice(&device)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}
			var deviceJson DeviceJsonReturn
			err = json.Unmarshal(j, &deviceJson)
			if err != nil {
				errs = errors.Join(err)
				continue
			}

			DbDevices, err := st.db.GetDevicesByHostname(device.Hostname)
			if err != nil {
				err = errors.New("error returning devices from db")
				errs = errors.Join(errs, err)
				continue
			}

			for _, dbDevice := range DbDevices {
				dbByte, err := st.translateDevice(&dbDevice)
				if err != nil {
					err = errors.New("error translating device")
					errs = errors.Join(errs, err)
					continue
				}
				var dbJson DeviceJsonReturn
				err = json.Unmarshal(dbByte, &dbJson)
				if err != nil {
					err = errors.New("error unmarshaling data")
					errs = errors.Join(errs, err)
					continue
				}
				ret, err := deviceJson.CheckDeviceEqual(deviceJson, dbJson)
				if err != nil {
					st.logger.Error("error checking device equality", zap.Any("error: ", err))
					continue
				}
				st.logger.Info("devices difference", zap.String("diffs: ", ret))
			}

			id, err := st.pusher.CreateDevice(j)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}

			// we store multiples device addresses because the agent can catch more than one per round
			// the ipChecker struct will be sent to the createInterfaceIpAddress as parameter
			ipChecker.IpInfo.DeviceAddresses = append(ipChecker.IpInfo.DeviceAddresses, deviceJson.IpAddress.Address) 
			ipChecker.IpInfo.DeviceId = append(ipChecker.IpInfo.DeviceId, id)
			newDevice, err := st.db.UpdateDevice(device.Id, id)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}
			if err := st.checkExistingInterfaces(&newDevice); err != nil {
				errs = errors.Join(errs, err)
				continue
			}
			if err := st.checkExistingInventories(&newDevice); err != nil {
				errs = errors.Join(errs, err)
				continue
			}

		}
		return errs
	} else if ifs, ok := data.([]storage.DbInterface); ok {

		var errs error
		for _, ifce := range ifs {
			if len(ifce.Id) == 0 {
				continue
			}
			device, err := st.db.GetDeviceByPolicyAndNamespaceAndHostname(ifce.Policy, ifce.Namespace, ifce.Hostname)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			} else if device.NetboxRefId == invalid_id {
				err = errors.New("invalid device id")
				errs = errors.Join(errs, err)
				continue
			}

			j, err := st.translateInterface(&ifce, device.NetboxRefId)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}

			DbInterfaces, err := st.db.GetInterfacesByName(ifce.Name)
			if err != nil {
				st.logger.Error("error retrieving interfaces", zap.Any("error: ", err))
				continue
			}

			for _, ifcDb := range DbInterfaces {
				ifceSq := ifce
				ifcSqTranslated, err := st.translateInterface(&ifceSq, device.NetboxRefId)
				if err != nil {
					st.logger.Error("error translating interface", zap.Any("error: ", err))
					continue
				}
				var ifceSqJson IfJsonReturn
				err = json.Unmarshal(ifcSqTranslated, &ifceSqJson)
				if err != nil {
					st.logger.Error("error unmarshaling interface", zap.Any("error: ", err))
					continue
				}

				ifcDbTranslated, err := st.translateInterface(&ifcDb, device.NetboxRefId)
				if err != nil {
					st.logger.Error("error translating interface", zap.Any("error: ", err))
					continue
				}
				var ifcDbJson IfJsonReturn
				err = json.Unmarshal(ifcDbTranslated, &ifcDbJson)
				if err != nil {
					st.logger.Error("error unmarshaling interface", zap.Any("error: ", err))
					continue
				}

				ret, err := ifcDbJson.CheckInterfaceEqual(ifcDbJson, ifceSqJson)
				if err != nil {
					st.logger.Error("error checking interface equality", zap.Any("error: ", err))
					continue
				}
				st.logger.Info("interfaces difference", zap.String("diffs: ", ret))
			}

			id, err := st.pusher.CreateInterface(j)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}
			newInterface, err := st.db.UpdateInterface(ifce.Id, id)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}
			for _, ip := range newInterface.IpAddresses {
				isMatch := false
				if device.Address == ip.Address {
					isMatch = true
				}

				j, err := st.translateIpInterface(&ip, newInterface.NetboxRefId, isMatch)
				if err != nil {
					errs = errors.Join(errs, err)
					continue
				}

				var ifIpJson IfIpJsonReturn
				err = json.Unmarshal(j, &ifIpJson)
				if err != nil {
					errs = errors.Join(errs, err)
				}
				_, err = st.pusher.CreateInterfaceIpAddress(j, ipChecker)
				if err != nil {
					errs = errors.Join(errs, err)
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
				errs = errors.Join(errs, err)
				continue
			}
		}
		return errs
	} else if inventories, ok := data.([]storage.DbInventory); ok {
		var errs error
		for _, inventory := range inventories {
			if len(inventory.Id) == 0 {
				continue
			}
			device, err := st.db.GetDeviceByPolicyAndNamespaceAndHostname(inventory.Policy, inventory.Namespace, inventory.Hostname)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			} else if device.NetboxRefId == invalid_id {
				err = errors.New("invalid device id")
				errs = errors.Join(errs, err)
				continue
			}
			j, err := st.translateInventory(&inventory, device.NetboxRefId)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}


			DbInventories, err := st.db.GetInventoriesByName(inventory.Name)
			if err != nil {
				st.logger.Error("error retrieving inventories", zap.Any("error: ", err))
				continue
			}

			for _, invDb := range DbInventories {
				invSq := inventory
				invSqTranslated, err := st.translateInventory(&invSq, device.NetboxRefId)
				if err != nil {
					st.logger.Error("error translating interface", zap.Any("error: ", err))
					continue
				}
				var invSqJson InvJsonReturn
				err = json.Unmarshal(invSqTranslated, &invSqJson)
				if err != nil {
					st.logger.Error("error unmarshaling interface", zap.Any("error: ", err))
					continue
				}

				invDbTranslated, err := st.translateInventory(&invDb, device.NetboxRefId)
				if err != nil {
					st.logger.Error("error translating interface", zap.Any("error: ", err))
					continue
				}
				var invDbJson InvJsonReturn
				err = json.Unmarshal(invDbTranslated, &invDbJson)
				if err != nil {
					st.logger.Error("error unmarshaling interface", zap.Any("error: ", err))
					continue
				}

				ret, err := invDbJson.CheckInventoriesEqual(invDbJson, invSqJson)
				if err != nil {
					st.logger.Error("error checking inventories equality", zap.Any("error: ", err))
					continue
				}
				st.logger.Info("devices difference", zap.String("diffs: ", ret))
			}

			id, err := st.pusher.CreateInventory(j)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}
			_, err = st.db.UpdateInventory(inventory.Id, id)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}
		}
		return errs
	}

	return errors.New("no valid translatable data found")
}

func (st *SuzieQTranslate) translateNetboxConfig(conf interface{}, reqKey string) interface{} {
	c, ok := conf.(map[string]interface{})
	if ok {
		n, ok := c["netbox"].(map[string]interface{})
		if ok {
			for k, v := range n {
				if k == reqKey {
					return v
				}
			}
		}
	}
	st.logger.Warn("value for the requered key not found", zap.String("key", reqKey))
	return nil
}

func (st *SuzieQTranslate) translateDevice(device *storage.DbDevice) ([]byte, error) {
	var ret DeviceJsonReturn
	if device.Config != nil {
		if value := st.translateNetboxConfig(device.Config, "site"); value != nil {
			if name, ok := value.(string); ok {
				ret.Site = &struct {
					Name string `json:"name"`
				}{Name: name}
			}
		}
	}
	ret.Name = device.Hostname
	ret.Status = device.State
	ret.Serial = device.SerialNumber
	ret.Dtype.Model = device.Model
	ret.Dtype.Mfr.Name = device.Vendor
	if len(device.Os) > 0 {
		ret.Platform = &struct {
			Name string `json:"name"`
			Mfr  struct {
				Name string `json:"name"`
			} `json:"manufacturers"`
		}{Name: device.Os}
		ret.Platform.Mfr.Name = device.Vendor
	}
	ret.IpAddress = &struct{
		Address string "json:\"ip_address,omitempty\""; 
		Version string "json:\"ip_version,omitempty\"";
		}{Address: device.Address, Version: ""}

	return json.Marshal(ret)
}

func (st *SuzieQTranslate) translateInterface(ifs *storage.DbInterface, deviceID int64) ([]byte, error) {
	var ret IfJsonReturn
	ret.DeviceID = deviceID
	ret.Name = ifs.Name
	ret.Type = ifs.IfType
	ret.Speed = ifs.Speed
	ret.Mtu = ifs.Mtu
	ret.MacAddress = ifs.MacAddress
	ret.State = ifs.AdminState
	return json.Marshal(ret)
}

func (st *SuzieQTranslate) translateIpInterface(ip *storage.IpAddress, ifID int64, isMatch bool) ([]byte, error) {
	var ret IfIpJsonReturn
	ret.IfID = ifID
	ret.Ip = ip.Address
	ret.IsPrimaryIp = isMatch
	return json.Marshal(ret)
}

func (st *SuzieQTranslate) translateVlan(vlan *storage.DbVlan) ([]byte, error) {
	st.logger.Warn("translate for vlan data not implemented yet")
	return nil, nil
}

func (st *SuzieQTranslate) translateInventory(inv *storage.DbInventory, deviceID int64) ([]byte, error) {
	var ret InvJsonReturn
	ret.DeviceID = deviceID
	ret.Name = inv.Name
	ret.AssetTag = inv.Type
	ret.Mfr.Name = inv.Vendor
	ret.Descr = inv.Descr
	ret.PartId = inv.PartNum
	ret.Serial = inv.Serial
	return json.Marshal(ret)
}

func (st *SuzieQTranslate) checkExistingInterfaces(device *storage.DbDevice) error {
	ifs, err := st.db.GetInterfaceByPolicyAndNamespaceAndHostname(device.Policy, device.Namespace, device.Hostname)
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

func (st *SuzieQTranslate) checkExistingInventories(device *storage.DbDevice) error {
	invs, err := st.db.GetInventoriesByPolicyAndNamespaceAndHostname(device.Policy, device.Namespace, device.Hostname)
	if err != nil {
		return nil
	}
	var inv []storage.DbInventory
	for _, v := range invs {
		if v.NetboxRefId == invalid_id {
			inv = append(inv, v)
		}
	}
	if len(inv) > 0 {
		return st.Translate(inv)
	}
	return nil
}
