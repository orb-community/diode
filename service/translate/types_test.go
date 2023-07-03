package translate

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckDeviceEqual(t *testing.T) {
	var sb strings.Builder
	databaseDevice := DeviceJsonReturn{
		Name:   "name_test",
		Status: "status_test",
		Serial: "serial_test",
		Dtype: struct {
			Model string "json:\"model\""
			Mfr   struct {
				Name string "json:\"name\""
			} "json:\"manufacturers\""
		}{
			Model: "model_test",
			Mfr: struct {
				Name string "json:\"name\""
			}{
				Name: "mfr_name_test",
			},
		},
		Platform: &struct {
			Name string "json:\"name\""
			Mfr  struct {
				Name string "json:\"name\""
			} "json:\"manufacturers\""
		}{
			Name: "platform_test",
			Mfr: struct {
				Name string "json:\"name\""
			}{
				Name: "mfr_name_test",
			},
		},
		Site: &struct {
			Name string "json:\"name\""
		}{
			Name: "site_name_test",
		},
	}
	suzieQdevice := DeviceJsonReturn{
		Name:   "name_test",
		Status: "status_test",
		Serial: "serial_test",
		Dtype: struct {
			Model string "json:\"model\""
			Mfr   struct {
				Name string "json:\"name\""
			} "json:\"manufacturers\""
		}{
			Model: "model_test",
			Mfr: struct {
				Name string "json:\"name\""
			}{
				Name: "mfr_name_test",
			},
		},
		Platform: &struct {
			Name string "json:\"name\""
			Mfr  struct {
				Name string "json:\"name\""
			} "json:\"manufacturers\""
		}{
			Name: "platform_test",
			Mfr: struct {
				Name string "json:\"name\""
			}{
				Name: "mfr_name_test",
			},
		},
		Site: &struct {
			Name string "json:\"name\""
		}{
			Name: "site_name_test",
		},
	}

	sb.WriteString("Devices are totally equal\n")

	got, err := databaseDevice.CheckDeviceEqual(suzieQdevice, databaseDevice)
	stringRet := sb.String()

	assert.Nil(t, err)
	assert.Equal(t, stringRet, got)
}

func TestCheckDeviceNotEqual(t *testing.T) {
	var sb strings.Builder
	dbDevice := DeviceJsonReturn{
		Name:   "name_test",
		Status: "status_test_another",
		Serial: "serial_test_another",
		Dtype: struct {
			Model string "json:\"model\""
			Mfr   struct {
				Name string "json:\"name\""
			} "json:\"manufacturers\""
		}{
			Model: "model_test_another",
			Mfr: struct {
				Name string "json:\"name\""
			}{
				Name: "mfr_name_test_another",
			},
		},
		Platform: &struct {
			Name string "json:\"name\""
			Mfr  struct {
				Name string "json:\"name\""
			} "json:\"manufacturers\""
		}{
			Name: "platform_test_another",
			Mfr: struct {
				Name string "json:\"name\""
			}{
				Name: "mfr_name_test_another",
			},
		},
		Site: &struct {
			Name string "json:\"name\""
		}{
			Name: "site_name_test_another",
		},
	}

	sqDevice := DeviceJsonReturn{
		Name:   "name_test",
		Status: "status_test",
		Serial: "serial_test",
		Dtype: struct {
			Model string "json:\"model\""
			Mfr   struct {
				Name string "json:\"name\""
			} "json:\"manufacturers\""
		}{
			Model: "model_test",
			Mfr: struct {
				Name string "json:\"name\""
			}{
				Name: "mfr_name_test",
			},
		},
		Platform: &struct {
			Name string "json:\"name\""
			Mfr  struct {
				Name string "json:\"name\""
			} "json:\"manufacturers\""
		}{
			Name: "platform_test",
			Mfr: struct {
				Name string "json:\"name\""
			}{
				Name: "mfr_name_test",
			},
		},
		Site: &struct {
			Name string "json:\"name\""
		}{
			Name: "site_name_test",
		},
	}
	var DiffJsonRet DiffJsonDeviceRet
	var DiffDev DiffsDevice

	DiffJsonRet.MfrName.NetBoxMfrName = dbDevice.Dtype.Mfr.Name
	DiffJsonRet.MfrName.SuzieQMfrName = sqDevice.Dtype.Mfr.Name
	DiffJsonRet.Model.NetBoxModel = dbDevice.Dtype.Model
	DiffJsonRet.Model.SuzieQModel = sqDevice.Dtype.Model
	DiffJsonRet.Platform.NetBoxPltName = dbDevice.Platform.Name
	DiffJsonRet.Platform.SuzieQPltName = sqDevice.Platform.Name
	DiffJsonRet.Status.NetBoxStatus = dbDevice.Status
	DiffJsonRet.Status.SuzieQStatus = sqDevice.Status
	DiffJsonRet.Serial.NetBoxSerial = dbDevice.Serial
	DiffJsonRet.Serial.SuzieQSerial = sqDevice.Serial
	DiffJsonRet.SiteName.NetBoxSiteName = dbDevice.Site.Name
	DiffJsonRet.SiteName.SuzieQSiteName = sqDevice.Site.Name
	DiffJsonRet.Name.NetBoxName = dbDevice.Name
	DiffJsonRet.Name.SuzieQName = sqDevice.Name

	DiffDev.DeviceDiffs = append(DiffDev.DeviceDiffs, DiffJsonRet)
	DiffDevJson, _ := json.Marshal(DiffDev)

	sb.WriteString(string(DiffDevJson))
	got, err := dbDevice.CheckDeviceEqual(sqDevice, dbDevice)
	stringRet := sb.String()

	assert.Nil(t, err)
	assert.Equal(t, stringRet, got)
}
func TestCheckDeviceWithEmptyObject(t *testing.T) {
	var sb strings.Builder
	dbDevice := DeviceJsonReturn{
		Name:   "name_test",
		Status: "status_test_another",
		Serial: "serial_test_another",
		Dtype: struct {
			Model string "json:\"model\""
			Mfr   struct {
				Name string "json:\"name\""
			} "json:\"manufacturers\""
		}{
			Model: "model_test_another",
			Mfr: struct {
				Name string "json:\"name\""
			}{
				Name: "mfr_name_test_another",
			},
		},
		Platform: &struct {
			Name string "json:\"name\""
			Mfr  struct {
				Name string "json:\"name\""
			} "json:\"manufacturers\""
		}{
			Name: "platform_test_another",
			Mfr: struct {
				Name string "json:\"name\""
			}{
				Name: "mfr_name_test_another",
			},
		},
		Site: &struct {
			Name string "json:\"name\""
		}{
			Name: "site_name_test_another",
		},
	}

	sqDevice := DeviceJsonReturn{}


	sb.WriteString("Empty object catched by Suzie Q")

	got, err := dbDevice.CheckDeviceEqual(sqDevice, dbDevice)
	stringRet := sb.String()

	assert.Nil(t, err)
	assert.Equal(t, stringRet, got)

}

func TestCheckInterfaceEqual(t *testing.T) {
	var sb strings.Builder
	databaseIfc := IfJsonReturn{
		DeviceID:   1,
		Name:       "test_name",
		Type:       "test_type",
		Speed:      200,
		Mtu:        300,
		MacAddress: "test_addr",
		State:      "test_state",
	}

	suzieQIfc := IfJsonReturn{
		DeviceID:   1,
		Name:       "test_name",
		Type:       "test_type",
		Speed:      200,
		Mtu:        300,
		MacAddress: "test_addr",
		State:      "test_state",
	}

	got, err := suzieQIfc.CheckInterfaceEqual(suzieQIfc, databaseIfc)

	sb.WriteString("Interfaces are totally equal\n")
	retString := sb.String()

	assert.Nil(t, err)
	assert.Equal(t, retString, got)
}

func TestCheckInterfaceNotEqual(t *testing.T) {
	var sb strings.Builder
	dbInterface := IfJsonReturn{
		DeviceID:   1000,
		Name:       "test_name",
		Type:       "test_type",
		Speed:      200,
		Mtu:        300,
		MacAddress: "test_addr",
		State:      "test_state",
	}

	sqInterface := IfJsonReturn{
		DeviceID:   1100,
		Name:       "test_name",
		Type:       "test_type_other",
		Speed:      300,
		Mtu:        500,
		MacAddress: "test_addr_other",
		State:      "test_state_other",
	}
	
	var ifcDiffs DiffInterfaceRet
	var IfcDiffsReturn DiffsInterface

	ifcDiffs.MacAddress.SuzieQIfcMacAddr = sqInterface.MacAddress
	ifcDiffs.MacAddress.NetBoxIfcMacAddr = dbInterface.MacAddress

	ifcDiffs.Speed.SuzieQIfcSpeed = sqInterface.Speed
	ifcDiffs.Speed.NetboxIfcSpeed = dbInterface.Speed

	ifcDiffs.DeviceID.SuzieQDevId = sqInterface.DeviceID
	ifcDiffs.DeviceID.NetBoxDevId = dbInterface.DeviceID

	ifcDiffs.Type.SuzieQIfcType = sqInterface.Type
	ifcDiffs.Type.NetBoxIfcType = dbInterface.Type

	ifcDiffs.Mtu.SuzieQIfcMtu = sqInterface.Mtu
	ifcDiffs.Mtu.NetBoxIfcMtu = dbInterface.Mtu

	ifcDiffs.State.NetBoxIfcState = dbInterface.State
	ifcDiffs.State.SuzieQIfcState = sqInterface.State

	IfcDiffsReturn.InterfaceDiffs = append(IfcDiffsReturn.InterfaceDiffs, ifcDiffs)

	DiffIfcJson, _ := json.Marshal(IfcDiffsReturn)

	sb.WriteString(string(DiffIfcJson))

	got, err := sqInterface.CheckInterfaceEqual(sqInterface, dbInterface)
	stringRet := sb.String()
	assert.Nil(t, err)
	assert.Equal(t, stringRet, got)
}

func TestCheckInterfaceWithEmptyObject(t *testing.T) {
	var sb strings.Builder
	databaseIfc := IfJsonReturn{
		DeviceID:   1,
		Name:       "test_name",
		Type:       "test_type",
		Speed:      200,
		Mtu:        300,
		MacAddress: "test_addr",
		State:      "test_state",
	}

	suzieQIfc := IfJsonReturn{}

	got, err := suzieQIfc.CheckInterfaceEqual(suzieQIfc, databaseIfc)

	sb.WriteString("Empty object catched by Suzie Q")
	retString := sb.String()

	assert.Nil(t, err)
	assert.Equal(t, retString, got)
}

func TestCheckInventoriesEqual(t *testing.T) {
	var sb strings.Builder
	sqInv := InvJsonReturn{
		DeviceID: 1,
		Name:     "name_test",
		AssetTag: "asset_test",
		Mfr: struct {
			Name string "json:\"name\""
		}{
			Name: "mfr_name",
		},
		PartId: "partid_test",
		Descr:  "descr_test",
		Serial: "serial_test",
	}

	dbInv := InvJsonReturn{
		DeviceID: 1,
		Name:     "name_test",
		AssetTag: "asset_test",
		Mfr: struct {
			Name string "json:\"name\""
		}{
			Name: "mfr_name",
		},
		PartId: "partid_test",
		Descr:  "descr_test",
		Serial: "serial_test",
	}

	sb.WriteString("Inventories are totally equal\n")

	got, err := sqInv.CheckInventoriesEqual(sqInv, dbInv)
	retString := sb.String()
	assert.Nil(t, err)
	assert.Equal(t, retString, got)
}

func TestCheckInventoriesNotEqual(t *testing.T) {
	var sb strings.Builder
	sqInv := InvJsonReturn{
		DeviceID: 1,
		Name:     "name_test",
		AssetTag: "asset_test",
		Mfr: struct {
			Name string "json:\"name\""
		}{
			Name: "mfr_name",
		},
		PartId: "partid_test",
		Descr:  "descr_test",
		Serial: "serial_test",
	}

	dbInv := InvJsonReturn{
		DeviceID: 2,
		Name:     "name_test",
		AssetTag: "asset_test_other",
		Mfr: struct {
			Name string "json:\"name\""
		}{
			Name: "mfr_name_other",
		},
		PartId: "partid_test_other",
		Descr:  "descr_test_other",
		Serial: "serial_test_other",
	}
  
	var DiffsRet DiffsInvRet

	var DiffsInv DiffInventoriesRet

	DiffsInv.DeviceID.SuzieQDevId = sqInv.DeviceID
	DiffsInv.DeviceID.NetBoxDevId = dbInv.DeviceID


	DiffsInv.AssetTag.SuzieQAssetTag = sqInv.AssetTag
	DiffsInv.AssetTag.NetBoxAssetTag = dbInv.AssetTag


	DiffsInv.MfrName.SuzieQMfrName = sqInv.Mfr.Name
	DiffsInv.MfrName.NetBoxMfrName = dbInv.Mfr.Name


	DiffsInv.PartId.SuzieQPartId = sqInv.PartId
	DiffsInv.PartId.NetBoxPartid = dbInv.PartId


	DiffsInv.Descr.SuzieQDescr = sqInv.Descr
	DiffsInv.Descr.NetBoxDescr = dbInv.Descr

	DiffsInv.Serial.SuzieQSerial = sqInv.Serial
	DiffsInv.Serial.NetBoxSerial = dbInv.Serial

	DiffsInv.Name.NetBoxIfcName = dbInv.Name
	DiffsInv.Name.SuzieQIfcName = sqInv.Name
	
	DiffsRet.InventoriesDiffs = append(DiffsRet.InventoriesDiffs, DiffsInv)

	DiffInvJson, _ := json.Marshal(DiffsRet)

	sb.WriteString(string(DiffInvJson))

	got, err := sqInv.CheckInventoriesEqual(sqInv, dbInv)
	retString := sb.String()
	assert.Nil(t, err)
	assert.Equal(t, retString, got)
}

func TestCheckInventoryWithEmptyObject(t *testing.T) {
	var sb strings.Builder
	sqInv := InvJsonReturn{
		DeviceID: 1,
		Name:     "name_test",
		AssetTag: "asset_test",
		Mfr: struct {
			Name string "json:\"name\""
		}{
			Name: "mfr_name",
		},
		PartId: "partid_test",
		Descr:  "descr_test",
		Serial: "serial_test",
	}

	dbInv := InvJsonReturn{}


	sb.WriteString("Empty object catched by Suzie Q")

	got, err := dbInv.CheckInventoriesEqual(dbInv, sqInv)
	stringRet := sb.String()

	assert.Nil(t, err)
	assert.Equal(t, stringRet, got)

}
