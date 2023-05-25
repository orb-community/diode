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
	dbJson, _ := json.Marshal(databaseDevice)
	sqJson, _ := json.Marshal(suzieQdevice)

	sb.WriteString("Database Device: ")
	sb.WriteString(string(dbJson) + "\n")
	sb.WriteString("Suzie Q catched Device: ")
	sb.WriteString(string(sqJson) + "\n")

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
		Serial: "serial_serial_another",
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
	sb.WriteString("DIFF: \n")
	sb.WriteString("Database Device Manufacture Name: ")
	sb.WriteString(string(sqDevice.Dtype.Mfr.Name) + "\n")
	sb.WriteString("Suzie Q Device Manufacture Name: ")
	sb.WriteString(string(dbDevice.Dtype.Mfr.Name) + "\n")

	sb.WriteString("Database Device Model Type: ")
	sb.WriteString(string(sqDevice.Dtype.Model) + "\n")
	sb.WriteString("Suzie Q Device Model Type: ")
	sb.WriteString(string(dbDevice.Dtype.Model) + "\n")

	sb.WriteString("Database Device Platform Name: ")
	sb.WriteString(string(sqDevice.Platform.Name) + "\n")
	sb.WriteString("Suzie Q Device Platform Name: ")
	sb.WriteString(string(dbDevice.Platform.Name) + "\n")

	sb.WriteString("Database Device Status: ")
	sb.WriteString(string(sqDevice.Status) + "\n")
	sb.WriteString("Suzie Q Device Status: ")
	sb.WriteString(string(dbDevice.Status) + "\n")

	sb.WriteString("Database Device Serial: ")
	sb.WriteString(string(sqDevice.Serial) + "\n")
	sb.WriteString("Suzie Q Device Serial: ")
	sb.WriteString(string(dbDevice.Serial) + "\n")

	sb.WriteString("Database Device Site Name: ")
	sb.WriteString(string(sqDevice.Site.Name) + "\n")
	sb.WriteString("Suzie Q Device Site Name: ")
	sb.WriteString(string(dbDevice.Site.Name) + "\n")
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
	dbJson, _ := json.Marshal(suzieQIfc)
	sqJson, _ := json.Marshal(databaseIfc)

	sb.WriteString("Database Interface: ")
	sb.WriteString(string(dbJson) + "\n")
	sb.WriteString("Suzie Q catched Interface: ")
	sb.WriteString(string(sqJson) + "\n")
	retString := sb.String()

	assert.Nil(t, err)
	assert.Equal(t, retString, got)
}

func TestCheckInterfaceNotEqual(t *testing.T) {
	var sb strings.Builder
	dbInterface := IfJsonReturn{
		DeviceID:   1,
		Name:       "test_name",
		Type:       "test_type",
		Speed:      200,
		Mtu:        300,
		MacAddress: "test_addr",
		State:      "test_state",
	}

	sqInterface := IfJsonReturn{
		DeviceID:   1,
		Name:       "test_name",
		Type:       "test_type_other",
		Speed:      200,
		Mtu:        300,
		MacAddress: "test_addr_other",
		State:      "test_state_other",
	}
	sb.WriteString("DIFF: \n")
	sb.WriteString("Database IFC Mac Address: ")
	sb.WriteString(string(dbInterface.MacAddress) + "\n")
	sb.WriteString("Suzie Q IFC Mac Address: ")
	sb.WriteString(string(sqInterface.MacAddress) + "\n")

	sb.WriteString("Database IFC State: ")
	sb.WriteString(string(dbInterface.State) + "\n")
	sb.WriteString("Suzie Q IFC State: ")
	sb.WriteString(string(sqInterface.State) + "\n")

	sb.WriteString("Database IFC Type: ")
	sb.WriteString(string(dbInterface.Type) + "\n")
	sb.WriteString("Suzie Q IFC Type: ")
	sb.WriteString(string(sqInterface.Type) + "\n")

	got, err := sqInterface.CheckInterfaceEqual(sqInterface, dbInterface)
	stringRet := sb.String()

	assert.Nil(t, err)
	assert.Equal(t, stringRet, got)
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
	dbJson, _ := json.Marshal(dbInv)
	sqJson, _ := json.Marshal(sqInv)

	sb.WriteString("Database Inventory: ")
	sb.WriteString(string(dbJson) + "\n")
	sb.WriteString("Suzie Q Inventory: ")
	sb.WriteString(string(sqJson) + "\n")

	got, err := sqInv.CheckInventoriesEqual(sqInv, dbInv)
	retString := sb.String()
	assert.Nil(t, err)
	assert.Equal(t, retString, got)
}

func TestCheckIventoriesNotEqual(t *testing.T) {
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
		AssetTag: "asset_test_other",
		Mfr: struct {
			Name string "json:\"name\""
		}{
			Name: "mfr_name",
		},
		PartId: "partid_test",
		Descr:  "descr_test_other",
		Serial: "serial_test_other",
	}

	sb.WriteString("DIFF: \n")

	sb.WriteString("Database INV AssetTag: ")
	sb.WriteString(string(dbInv.AssetTag) + "\n")
	sb.WriteString("Suzie Q INV AssetTag: ")
	sb.WriteString(string(sqInv.AssetTag) + "\n")
	sb.WriteString("Database INV Description: ")
	sb.WriteString(string(dbInv.Descr) + "\n")
	sb.WriteString("Suzie Q INV Description: ")
	sb.WriteString(string(sqInv.Descr) + "\n")
	sb.WriteString("Database INV Serial: ")
	sb.WriteString(string(dbInv.Serial) + "\n")
	sb.WriteString("Suzie Q INV Serial: ")
	sb.WriteString(string(sqInv.Serial) + "\n")

	got, err := sqInv.CheckInventoriesEqual(sqInv, dbInv)
	retString := sb.String()
	assert.Nil(t, err)
	assert.Equal(t, retString, got)

}
