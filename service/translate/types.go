package translate

import (
	"encoding/json"
	"fmt"
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

func (DeviceJsonReturn) CheckDeviceEqual(sqDevice, dbDevice DeviceJsonReturn) (string, error) {
	var sb strings.Builder
	if dbDevice.Name == sqDevice.Name {
		if reflect.DeepEqual(sqDevice, dbDevice) {
			sb.WriteString("Devices are totally equal\n")
			dbJson, err := json.Marshal(dbDevice)
			if err != nil {
				return "", err
			}
			sqJson, err := json.Marshal(sqDevice)
			if err != nil {
				return "", err
			}
			sb.WriteString("Database Device: ")
			sb.WriteString(string(dbJson) + "\n")
			sb.WriteString("Suzie Q catched Device: ")
			sb.WriteString(string(sqJson) + "\n")
			return sb.String(), nil
		} else {
			sb.WriteString("DIFF: \n")
			if sqDevice.Dtype.Mfr.Name != dbDevice.Dtype.Mfr.Name {
				sb.WriteString("Database Device Manufacture Name: ")
				sb.WriteString(string(sqDevice.Dtype.Mfr.Name) + "\n")
				sb.WriteString("Suzie Q Device Manufacture Name: ")
				sb.WriteString(string(dbDevice.Dtype.Mfr.Name) + "\n")
			}
			if sqDevice.Dtype.Model != dbDevice.Dtype.Model {
				sb.WriteString("Database Device Model Type: ")
				sb.WriteString(string(sqDevice.Dtype.Model) + "\n")
				sb.WriteString("Suzie Q Device Model Type: ")
				sb.WriteString(string(dbDevice.Dtype.Model) + "\n")
			}
			if sqDevice.Platform.Name != dbDevice.Platform.Name {
				sb.WriteString("Database Device Platform Name: ")
				sb.WriteString(string(sqDevice.Platform.Name) + "\n")
				sb.WriteString("Suzie Q Device Platform Name: ")
				sb.WriteString(string(dbDevice.Platform.Name) + "\n")
			}
			if sqDevice.Status != dbDevice.Status {
				sb.WriteString("Database Device Status: ")
				sb.WriteString(string(sqDevice.Status) + "\n")
				sb.WriteString("Suzie Q Device Status: ")
				sb.WriteString(string(dbDevice.Status) + "\n")
			}
			if sqDevice.Serial != dbDevice.Serial {
				sb.WriteString("Database Device Serial: ")
				sb.WriteString(string(sqDevice.Serial) + "\n")
				sb.WriteString("Suzie Q Device Serial: ")
				sb.WriteString(string(dbDevice.Serial) + "\n")
			}
			if sqDevice.Site.Name != dbDevice.Site.Name {
				sb.WriteString("Database Device Site Name: ")
				sb.WriteString(string(sqDevice.Site.Name) + "\n")
				sb.WriteString("Suzie Q Device Site Name: ")
				sb.WriteString(string(dbDevice.Site.Name) + "\n")
			}
			return sb.String(), nil
		}
	}
	return "", nil
}

func (IfJsonReturn) CheckInterfaceEqual(sqInterface, dbInterface IfJsonReturn) (string, error) {
	var sb strings.Builder
	if sqInterface.Name == dbInterface.Name {
		if reflect.DeepEqual(sqInterface, dbInterface) {
			sb.WriteString("Interfaces are totally equal\n")
			dbJson, err := json.Marshal(sqInterface)
			if err != nil {
				return "", err
			}
			sqJson, err := json.Marshal(dbInterface)
			if err != nil {
				return "", err
			}
			sb.WriteString("Database Interface: ")
			sb.WriteString(string(dbJson) + "\n")
			sb.WriteString("Suzie Q catched Interface: ")
			sb.WriteString(string(sqJson) + "\n")
			return sb.String(), nil
		} else {
			sb.WriteString("DIFF: \n")
			if sqInterface.MacAddress != dbInterface.MacAddress {
				sb.WriteString("Database IFC Mac Address: ")
				sb.WriteString(string(dbInterface.MacAddress) + "\n")
				sb.WriteString("Suzie Q IFC Mac Address: ")
				sb.WriteString(string(sqInterface.MacAddress) + "\n")
			}
			if sqInterface.Speed != dbInterface.Speed {
				sb.WriteString("Database IFC Speed: ")
				sb.WriteString(string(rune(dbInterface.Speed)) + "\n")
				sb.WriteString("Suzie Q IFC Speed: ")
				sb.WriteString(string(rune(sqInterface.Speed)) + "\n")
			}
			if sqInterface.DeviceID != dbInterface.DeviceID {
				sb.WriteString("Database IFC DeviceID: ")
				sb.WriteString(string(rune(dbInterface.DeviceID)) + "\n")
				sb.WriteString("Suzie Q IFC DeviceID: ")
				sb.WriteString(string(rune(sqInterface.DeviceID)) + "\n")
			}
			if sqInterface.State != dbInterface.State {
				sb.WriteString("Database IFC State: ")
				sb.WriteString(string(dbInterface.State) + "\n")
				sb.WriteString("Suzie Q IFC State: ")
				sb.WriteString(string(sqInterface.State) + "\n")
			}
			if sqInterface.Type != dbInterface.Type {
				sb.WriteString("Database IFC Type: ")
				sb.WriteString(string(dbInterface.Type) + "\n")
				sb.WriteString("Suzie Q IFC Type: ")
				sb.WriteString(string(sqInterface.Type) + "\n")
			}
			if sqInterface.Mtu != dbInterface.Mtu {
				sb.WriteString("Database IFC Mtu: ")
				sb.WriteString(string(rune(dbInterface.Mtu)) + "\n")
				sb.WriteString("Suzie Q IFC Mtu: ")
				sb.WriteString(string(rune(sqInterface.Mtu)) + "\n")
			}
			return sb.String(), nil
		}
	}
	return "", nil
}

func (InvJsonReturn) CheckInventoriesEqual(sqInventory, dbInventory InvJsonReturn) (string, error) {
	var sb strings.Builder
	if sqInventory.Name == dbInventory.Name {
		if reflect.DeepEqual(sqInventory, dbInventory) {
			sb.WriteString("Inventories are totally equal\n")
			dbJson, err := json.Marshal(dbInventory)
			if err != nil {
				return "", err
			}
			sqJson, err := json.Marshal(sqInventory)
			if err != nil {
				return "", err
			}
			sb.WriteString("Database Inventory: ")
			sb.WriteString(string(dbJson) + "\n")
			sb.WriteString("Suzie Q Inventory: ")
			sb.WriteString(string(sqJson) + "\n")
			fmt.Println(sb.String())
			return sb.String(), nil
		} else {
			sb.WriteString("DIFF: \n")
			if sqInventory.DeviceID != dbInventory.DeviceID {
				sb.WriteString("Database INV DeviceID: ")
				sb.WriteString(string(rune(dbInventory.DeviceID)) + "\n")
				sb.WriteString("Suzie Q INV DeviceID: ")
				sb.WriteString(string(rune(sqInventory.DeviceID)) + "\n")
			}
			if sqInventory.AssetTag != dbInventory.AssetTag {
				sb.WriteString("Database INV AssetTag: ")
				sb.WriteString(string(dbInventory.AssetTag) + "\n")
				sb.WriteString("Suzie Q INV AssetTag: ")
				sb.WriteString(string(sqInventory.AssetTag) + "\n")
			}
			if sqInventory.Mfr.Name != dbInventory.Mfr.Name {
				sb.WriteString("Database INV Mfr Name: ")
				sb.WriteString(string(dbInventory.Mfr.Name) + "\n")
				sb.WriteString("Suzie Q INV Mfr Name: ")
				sb.WriteString(string(sqInventory.Mfr.Name) + "\n")
			}
			if sqInventory.PartId != dbInventory.PartId {
				sb.WriteString("Database INV PartId: ")
				sb.WriteString(string(dbInventory.PartId) + "\n")
				sb.WriteString("Suzie Q INV PartId: ")
				sb.WriteString(string(sqInventory.PartId) + "\n")
			}
			if sqInventory.Descr != dbInventory.Descr {
				sb.WriteString("Database INV Description: ")
				sb.WriteString(string(dbInventory.Descr) + "\n")
				sb.WriteString("Suzie Q INV Description: ")
				sb.WriteString(string(sqInventory.Descr) + "\n")
			}
			if sqInventory.Serial != dbInventory.Serial {
				sb.WriteString("Database INV Serial: ")
				sb.WriteString(string(dbInventory.Serial) + "\n")
				sb.WriteString("Suzie Q INV Serial: ")
				sb.WriteString(string(sqInventory.Serial) + "\n")
			}
			return sb.String(), nil
		}
	}
	return "", nil
}
