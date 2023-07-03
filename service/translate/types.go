package translate

import (
	"encoding/json"
	"reflect"
	"strings"
)

type DeviceJsonReturn struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Serial string `json:"serial"`
	Dtype  struct {
		Model string `json:"model"`
		Mfr   struct {
			Name string `json:"name"`
		} `json:"manufacturers"`
	} `json:"device_type"`
	Platform *struct {
		Name string `json:"name"`
		Mfr  struct {
			Name string `json:"name"`
		} `json:"manufacturers"`
	} `json:"platform,omitempty"`
	Site *struct {
		Name string `json:"name"`
	} `json:"site,omitempty"`
	IpAddress *struct {
		Address string `json:"ip_address,omitempty"`
		Version string `json:"ip_version,omitempty"`
	} `json:"primary_ip,omitempty"`
 
}

type IfJsonReturn struct {
	DeviceID   int64  `json:"device_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Speed      int64  `json:"speed"`
	Mtu        int64  `json:"mtu"`
	MacAddress string `json:"mac_address"`
	State      string `json:"state"`
}

type IfIpJsonReturn struct {
	IfID int64  `json:"assigned_object_id"`
	Ip   string `json:"address"`
}

type InvJsonReturn struct {
	DeviceID int64  `json:"device_id"`
	Name     string `json:"name"`
	AssetTag string `json:"asset_tag"`
	Mfr      struct {
		Name string `json:"name"`
	} `json:"manufacturers"`
	PartId string `json:"part_id"`
	Descr  string `json:"description"`
	Serial string `json:"serial"`
}

type DiffsDevice struct {
	DeviceDiffs []DiffJsonDeviceRet `json:"device_diffs"`
}
type DiffJsonDeviceRet struct {
	Name struct {
		SuzieQName string `json:"suzieq_name,omitempty"`
		NetBoxName string `json:"netbox_name,omitempty"`
	} `json:"name,omitempty"`
	Status struct {
		SuzieQStatus string `json:"suzieq_status,omitempty"`
		NetBoxStatus string `json:"netbox_status,omitempty"`
	} `json:"dev_status,omitempty"`
	Serial struct {
		SuzieQSerial string `json:"suzieq_serial,omitempty"`
		NetBoxSerial string `json:"netbox_serial,omitempty"`
	} `json:"dev_serial,omitempty"`
	Model struct {
		SuzieQModel string `json:"suzieq_model,omitempty"`
		NetBoxModel string `json:"netbox_model,omitempty"`
	} `json:"dev_model,omitempty"`
	MfrName struct {
		SuzieQMfrName string `json:"suzieq_mfr_name,omitempty"`
		NetBoxMfrName string `json:"netbox_mfr_name,omitempty"`
	} `json:"dev_mfr_name,omitempty"`
	Platform struct {
		SuzieQPltName string `json:"suzieq_plt_name,omitempty"`
		NetBoxPltName string `json:"netbox_plt_name,omitempty"`
	} `json:"dev_plt,omitempty"`
	SiteName struct {
		SuzieQSiteName string `json:"suzieq_site_name,omitempty"`
		NetBoxSiteName string `json:"netbox_site_name,omitempty"`
	} `json:"dev_site_name,omitempty"`
	PrimaryIP struct {
		SuzieQAddress string `json:"suzieq_address,omitempty"`
		NetBoxAddress string `json:"netbox_address,omitempty"`
	} `json:"dev_primary_ip,omitempty"`
	Version struct {
		SuzieQVersion string `json:"suzieq_version,omitempty"`
		NetBoxVersion string `json:"netbox_version,omitempty"`
	} `json:"dev_version,omitempty"`
}

type DiffsInterface struct {
	InterfaceDiffs []DiffInterfaceRet `json:"interface_diffs"`
}

type DiffInterfaceRet struct {
	DeviceID struct {
		SuzieQDevId int64 `json:"suzieq_dev_id,omitempty"`
		NetBoxDevId int64 `json:"netbox_dev_id,omitempty"`
	} `json:"device_id,omitempty"`
	Type struct {
		SuzieQIfcType string `json:"suzieq_ifc_type,omitempty"`
		NetBoxIfcType string `json:"netbox_ifc_type,omitempty"`
	} `json:"ifc_type,omitempty"`
	Speed struct {
		SuzieQIfcSpeed int64 `json:"suzieq_ifc_speed,omitempty"`
		NetboxIfcSpeed int64 `json:"netbox_ifc_speed,omitempty"`
	} `json:"ifc_speed,omitempty"`
	Mtu struct {
		SuzieQIfcMtu int64 `json:"suzieq_ifc_mtu,omitempty"`
		NetBoxIfcMtu int64 `json:"netbox_ifc_mtu,omitempty"`
	} `json:"ifc_mtu,omitempty"`
	MacAddress struct {
		SuzieQIfcMacAddr string `json:"suzieq_ifc_mac_addr,omitempty"`
		NetBoxIfcMacAddr string `json:"netbox_ifc_mac_addr,omitempty"`
	} `json:"ifc_mac_addr,omitempty"`
	State struct {
		SuzieQIfcState string `json:"suzieq_ifc_state,omitempty"`
		NetBoxIfcState string `json:"netbox_ifc_state,omitempty"`
	} `json:"ifc_state,omitempty"`
}

type DiffsInvRet struct {
	InventoriesDiffs []DiffInventoriesRet `json:"inventories_diffs"`
}

type DiffInventoriesRet struct {
	DeviceID struct {
		SuzieQDevId int64 `json:"suzieq_dev_id,omitempty"`
		NetBoxDevId int64 `json:"netbox_dev_id,omitempty"`
	} `json:"device_id,omitempty"`
	Name struct {
		SuzieQIfcName string `json:"suzieq_inv_name,omitempty"`
		NetBoxIfcName string `json:"netbox_inv_name,omitempty"`
	} `json:"inv_name,omitempty"`
	AssetTag struct {
		SuzieQAssetTag string `json:"suzieq_asset_tag,omitempty"`
		NetBoxAssetTag string `json:"netbox_asset_tag,omitempty"`
	} `json:"asset_tag,omitempty"`
	MfrName struct {
		SuzieQMfrName string `json:"suzieq_mfr_name,omitempty"`
		NetBoxMfrName string `json:"netbox_mfr_name,omitempty"`
	} `json:"dev_mfr_name,omitempty"`
	PartId struct {
		SuzieQPartId string `json:"suzieq_partid,omitempty"`
		NetBoxPartid string `json:"netbox_partid,omitempty"`
	} `json:"part_id,omitempty"`
	Descr  struct {
		SuzieQDescr string `json:"suzieq_descr,omitempty"`
		NetBoxDescr string `json:"netbox_descr,omitempty"`
	} `json:"descr,omitempty"`
	Serial struct {
		SuzieQSerial string `json:"suzieq_serial,omitempty"`
		NetBoxSerial string `json:"netbox_serial,omitempty"`
	} `json:"dev_serial,omitempty"`
}

func (DeviceJsonReturn) CheckDeviceEqual(sqDevice, dbDevice DeviceJsonReturn) (string, error) {
	var sb strings.Builder
	if dbDevice.Name == sqDevice.Name {
		if reflect.DeepEqual(sqDevice, dbDevice) {
			sb.WriteString("Devices are totally equal\n")
			return sb.String(), nil
		} else {
			var DiffDev DiffsDevice
			var DiffJsonRet DiffJsonDeviceRet
			if sqDevice.Dtype.Mfr.Name != dbDevice.Dtype.Mfr.Name {
				DiffJsonRet.MfrName.NetBoxMfrName = dbDevice.Dtype.Mfr.Name
				DiffJsonRet.MfrName.SuzieQMfrName = sqDevice.Dtype.Mfr.Name
			}
			if sqDevice.Dtype.Model != dbDevice.Dtype.Model {
				DiffJsonRet.Model.NetBoxModel = dbDevice.Dtype.Model
				DiffJsonRet.Model.SuzieQModel = sqDevice.Dtype.Model
			}

			if sqDevice.Platform.Name != dbDevice.Platform.Name {
				DiffJsonRet.Platform.NetBoxPltName = dbDevice.Platform.Name
				DiffJsonRet.Platform.SuzieQPltName = sqDevice.Platform.Name
			}
			
			if sqDevice.Status != dbDevice.Status {
				DiffJsonRet.Status.NetBoxStatus = dbDevice.Status
				DiffJsonRet.Status.SuzieQStatus = sqDevice.Status
			}
			if sqDevice.Serial != dbDevice.Serial {
				DiffJsonRet.Serial.NetBoxSerial = dbDevice.Serial
				DiffJsonRet.Serial.SuzieQSerial = sqDevice.Serial
			}
			if sqDevice.IpAddress.Address != dbDevice.IpAddress.Address {
				DiffJsonRet.PrimaryIP.NetBoxAddress = dbDevice.IpAddress.Address
				DiffJsonRet.PrimaryIP.SuzieQAddress = sqDevice.IpAddress.Address
			}
			if sqDevice.IpAddress.Version != dbDevice.IpAddress.Version {
				DiffJsonRet.Version.NetBoxVersion = dbDevice.IpAddress.Version
				DiffJsonRet.Version.SuzieQVersion = sqDevice.IpAddress.Version
			}
			
			DiffJsonRet.Name.NetBoxName = sqDevice.Name
			DiffJsonRet.Name.SuzieQName = sqDevice.Name
			DiffDev.DeviceDiffs = append(DiffDev.DeviceDiffs, DiffJsonRet)
			DiffDevJson, err := json.Marshal(DiffDev)
			if err != nil {
				return "", err
			}
			sb.WriteString(string(DiffDevJson))

			return sb.String(), nil
		}
	} else if (DeviceJsonReturn{}) == sqDevice {
		sb.WriteString("Empty object catched by Suzie Q")
		return sb.String(), nil
	}
	return "", nil
}

func (IfJsonReturn) CheckInterfaceEqual(sqInterface, dbInterface IfJsonReturn) (string, error) {
	var sb strings.Builder
	if sqInterface.Name == dbInterface.Name {
		if reflect.DeepEqual(sqInterface, dbInterface) {
			sb.WriteString("Interfaces are totally equal\n")
			return sb.String(), nil
		} else {
			var ifcDiffs DiffInterfaceRet
			var IfcDiffsReturn DiffsInterface
			if sqInterface.MacAddress != dbInterface.MacAddress {
				ifcDiffs.MacAddress.SuzieQIfcMacAddr = sqInterface.MacAddress
				ifcDiffs.MacAddress.NetBoxIfcMacAddr = dbInterface.MacAddress
			}
			if sqInterface.Speed != dbInterface.Speed {
				ifcDiffs.Speed.SuzieQIfcSpeed = sqInterface.Speed
				ifcDiffs.Speed.NetboxIfcSpeed = dbInterface.Speed
			}
			if sqInterface.State != dbInterface.State {
				ifcDiffs.State.SuzieQIfcState = sqInterface.State
				ifcDiffs.State.NetBoxIfcState = dbInterface.State
			}
			if sqInterface.Type != dbInterface.Type {
				ifcDiffs.Type.SuzieQIfcType = sqInterface.Type
				ifcDiffs.Type.NetBoxIfcType = dbInterface.Type
			}
			if sqInterface.Mtu != dbInterface.Mtu {
				ifcDiffs.Mtu.SuzieQIfcMtu = sqInterface.Mtu
				ifcDiffs.Mtu.NetBoxIfcMtu = dbInterface.Mtu
			}
			
			IfcDiffsReturn.InterfaceDiffs = append(IfcDiffsReturn.InterfaceDiffs, ifcDiffs)

			DiffIfcJson, err := json.Marshal(IfcDiffsReturn)
			if err != nil {
				return "", err
			}
			sb.WriteString(string(DiffIfcJson))
			return sb.String(), nil
		}
	} else if (IfJsonReturn{}) == sqInterface {
		sb.WriteString("Empty object catched by Suzie Q")
		return sb.String(), nil
	}
	return "", nil
}

func (InvJsonReturn) CheckInventoriesEqual(sqInventory, dbInventory InvJsonReturn) (string, error) {
	var sb strings.Builder
	if sqInventory.Name == dbInventory.Name {
		if reflect.DeepEqual(sqInventory, dbInventory) {
			sb.WriteString("Inventories are totally equal\n")
			return sb.String(), nil
		} else {
			var DiffsRet DiffsInvRet
			var DiffsInv DiffInventoriesRet
			if sqInventory.AssetTag != dbInventory.AssetTag {
				DiffsInv.AssetTag.SuzieQAssetTag = sqInventory.AssetTag
				DiffsInv.AssetTag.NetBoxAssetTag = dbInventory.AssetTag
			}
			if sqInventory.Mfr.Name != dbInventory.Mfr.Name {
				DiffsInv.MfrName.SuzieQMfrName = sqInventory.Mfr.Name
				DiffsInv.MfrName.NetBoxMfrName = dbInventory.Mfr.Name
			}
			if sqInventory.PartId != dbInventory.PartId {
				DiffsInv.PartId.SuzieQPartId = sqInventory.PartId
				DiffsInv.PartId.NetBoxPartid = dbInventory.PartId
			}
			if sqInventory.Descr != dbInventory.Descr {
				DiffsInv.Descr.SuzieQDescr = sqInventory.Descr
				DiffsInv.Descr.NetBoxDescr = dbInventory.Descr
			}
			if sqInventory.Serial != dbInventory.Serial {
				DiffsInv.Serial.SuzieQSerial = sqInventory.Serial
				DiffsInv.Serial.NetBoxSerial = dbInventory.Serial
			}
			DiffsInv.Name.NetBoxIfcName = dbInventory.Name
			DiffsInv.Name.SuzieQIfcName = sqInventory.Name

			DiffsRet.InventoriesDiffs = append(DiffsRet.InventoriesDiffs, DiffsInv)

			DiffInvJson, err := json.Marshal(DiffsRet)
			if err != nil {
				return "", err
			}
			sb.WriteString(string(DiffInvJson))
			return sb.String(), nil
		}
	} else if (InvJsonReturn{}) == sqInventory {
		sb.WriteString("Empty object catched by Suzie Q")
		return sb.String(), nil
	}
	return "", nil
}
