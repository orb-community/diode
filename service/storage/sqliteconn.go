package storage

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

type sqliteStorage struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewSqliteStorage(logger *zap.Logger) (Service, error) {
	db, err := startSqliteDb(logger)
	if err != nil {
		return nil, err
	}
	return sqliteStorage{db: db, logger: logger}, nil
}



func (s sqliteStorage) GetInterfacesByName(name string) ([]DbInterface, error) {
	selectResult, err := s.db.Query(`
		SELECT id, policy, config, namespace, hostname, name, admin_state, mtu, speed, mac_address, if_type, netbox_id, ip_addresses, json_data
		FROM interfaces
		WHERE name = $1
	`, name)
	if err != nil {
		return nil, errors.Join(errors.New("storage fetch interface fail on new function"), err)
	}
	var interfaces []DbInterface
	var configAsString string
	var ipsAsString string
	for selectResult.Next() {
		var iface DbInterface
		err := selectResult.Scan(&iface.Id, &iface.Policy, &configAsString, &iface.Namespace, &iface.Hostname, &iface.Name, &iface.AdminState,
			&iface.Mtu, &iface.Speed, &iface.MacAddress, &iface.IfType, &iface.NetboxRefId, &ipsAsString, &iface.Blob)
		if err != nil {
			return nil, errors.Join(errors.New("storage ifce struct fail"), err)
		}
		if len(configAsString) > 0 {
			err = json.Unmarshal([]byte(configAsString), &iface.Config)
			if err != nil {
				return nil, errors.Join(errors.New("storage config parse fail"), err)
			}
		}
		err = json.Unmarshal([]byte(ipsAsString), &iface.IpAddresses)
		if err != nil {
			return nil, errors.Join(errors.New("storage ip_address parse fail"), err)
		}
		interfaces = append(interfaces, iface)
	}
	return interfaces, nil
}

func (s sqliteStorage) GetInterfaceByPolicyAndNamespaceAndHostname(policy, namespace, hostname string) ([]DbInterface, error) {
	selectResult, err := s.db.Query(`
		SELECT id, policy, config, namespace, hostname, name, admin_state, mtu, speed, mac_address, if_type, netbox_id, ip_addresses, json_data
		FROM interfaces
		WHERE policy = $1 AND namespace = $2 AND hostname = $3
	`, policy, namespace, hostname)
	if err != nil {
		return nil, errors.Join(errors.New("storage fetch interface fail"), err)
	}
	var interfaces []DbInterface
	var configAsString string
	var ipsAsString string
	for selectResult.Next() {
		var iface DbInterface
		err := selectResult.Scan(&iface.Id, &iface.Policy, &configAsString, &iface.Namespace, &iface.Hostname, &iface.Name, &iface.AdminState,
			&iface.Mtu, &iface.Speed, &iface.MacAddress, &iface.IfType, &iface.NetboxRefId, &ipsAsString, &iface.Blob)
		if err != nil {
			return nil, errors.Join(errors.New("storage ifce struct fail"), err)
		}
		if len(configAsString) > 0 {
			err = json.Unmarshal([]byte(configAsString), &iface.Config)
			if err != nil {
				return nil, errors.Join(errors.New("storage config parse fail"), err)
			}
		}
		err = json.Unmarshal([]byte(ipsAsString), &iface.IpAddresses)
		if err != nil {
			return nil, errors.Join(errors.New("storage ip_address parse fail"), err)
		}
		interfaces = append(interfaces, iface)
	}

	return interfaces, nil
}

func (s sqliteStorage) GetDevicesByHostname(hostname string) ([]DbDevice, error) {
	selectResult, err := s.db.Query(`
	SELECT id, policy, config, namespace, hostname, serial_number, model, state, vendor, os, netbox_id, json_data
	FROM devices
	WHERE hostname = $1`, hostname)
	if err != nil {
		return nil, errors.Join(errors.New("storage fetch device fail, not able to return devices"), err)
	}
	var devices []DbDevice
	var configAsString string
	for selectResult.Next() {
		var device DbDevice
		err := selectResult.Scan(&device.Id, &device.Policy, &configAsString, &device.Namespace, &device.Hostname, &device.SerialNumber,
			&device.Model, &device.State, &device.Vendor, &device.Os, &device.NetboxRefId, &device.Blob)
		if err != nil {
			return nil, errors.Join(errors.New("storage create device struct fail on NEW REPO FUNC"), err)
		}
		if len(configAsString) > 0 {
			err = json.Unmarshal([]byte(configAsString), &device.Config)
			if err != nil {
				return nil, errors.Join(errors.New("storage config parse fail NEW FUNC"), err)
			}
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (s sqliteStorage) GetDevicesByPolicyAndNamespace(policy, namespace string) ([]DbDevice, error) {
	selectResult, err := s.db.Query(`
		SELECT id, policy, config, namespace, hostname, serial_number, model, state, vendor, os, netbox_id, json_data
		FROM devices
		WHERE policy = $1 AND namespace = $2
	`, policy, namespace)
	if err != nil {
		return nil, errors.Join(errors.New("storage fetch device fail"), err)
	}
	var devices []DbDevice
	var configAsString string
	for selectResult.Next() {
		var device DbDevice
		err := selectResult.Scan(&device.Id, &device.Policy, &configAsString, &device.Namespace, &device.Hostname, &device.SerialNumber,
			&device.Model, &device.State, &device.Vendor, &device.Os, &device.NetboxRefId, &device.Blob)
		if err != nil {
			return nil, errors.Join(errors.New("storage create device struct fail"), err)
		}
		if len(configAsString) > 0 {
			err = json.Unmarshal([]byte(configAsString), &device.Config)
			if err != nil {
				return nil, errors.Join(errors.New("storage config parse fail"), err)
			}
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (s sqliteStorage) GetDeviceByPolicyAndNamespaceAndHostname(policy, namespace, hostname string) (DbDevice, error) {
	selectResult := s.db.QueryRow(`
		SELECT id, policy, config, namespace, hostname, serial_number, model, state, vendor, os, netbox_id, json_data
		FROM devices
		WHERE policy = $1 AND namespace = $2 AND hostname = $3
	`, policy, namespace, hostname)
	var device DbDevice
	var configAsString string
	err := selectResult.Scan(&device.Id, &device.Policy, &configAsString, &device.Namespace, &device.Hostname, &device.SerialNumber,
		&device.Model, &device.State, &device.Vendor, &device.Os, &device.NetboxRefId, &device.Blob)
	if err != nil {
		return DbDevice{}, errors.Join(errors.New("storage fetch device fail"), err)
	}
	if len(configAsString) > 0 {
		err = json.Unmarshal([]byte(configAsString), &device.Config)
		if err != nil {
			return DbDevice{}, errors.Join(errors.New("storage config parse fail"), err)
		}
	}
	return device, nil
}

func (s sqliteStorage) GetVlansByPolicyAndNamespaceAndHostname(policy, namespace, hostname string) ([]DbVlan, error) {
	selectResult, err := s.db.Query(`
		SELECT id, policy, config, namespace, hostname, name, state, netbox_id, json_data
		FROM vlans
		WHERE policy = $1 AND namespace = $2 AND hostname = $3
	`, policy, namespace, hostname)
	if err != nil {
		return nil, errors.Join(errors.New("storage fetch vlan fail"), err)
	}
	var vlans []DbVlan
	var configAsString string
	for selectResult.Next() {
		var vlan DbVlan
		err := selectResult.Scan(&vlan.Id, &vlan.Policy, &configAsString, &vlan.Namespace, &vlan.Hostname, &vlan.Name,
			&vlan.State, &vlan.NetboxRefId, &vlan.Blob)
		if err != nil {
			return nil, errors.Join(errors.New("storage create vlan struct fail"), err)
		}
		if len(configAsString) > 0 {
			err = json.Unmarshal([]byte(configAsString), &vlan.Config)
			if err != nil {
				return nil, errors.Join(errors.New("storage config parse fail"), err)
			}
		}
		vlans = append(vlans, vlan)
	}

	return vlans, nil
}

func (s sqliteStorage) GetInventoriesByName(name string) ([]DbInventory, error) {
	selectResult, err := s.db.Query(`
	SELECT id, policy, config, namespace, hostname, name, description, vendor, serial, part_num, type, netbox_id, json_data
	FROM inventories
	WHERE name = $1
	`, name)
	if err != nil {
		return nil, errors.Join(errors.New("storage fetch inventory fail"), err)
	}
	var inventories []DbInventory
	var configAsString string
	for selectResult.Next() {
		var inventory DbInventory
		err := selectResult.Scan(&inventory.Id, &inventory.Policy, &configAsString, &inventory.Namespace, &inventory.Hostname, &inventory.Name,
			&inventory.Descr, &inventory.Vendor, &inventory.Serial, &inventory.PartNum, &inventory.Type, &inventory.NetboxRefId, &inventory.Blob)
		if err != nil {
			return nil, errors.Join(errors.New("storage create inventory struct fail"), err)
		}
		if len(configAsString) > 0 {
			err = json.Unmarshal([]byte(configAsString), &inventory.Config)
			if err != nil {
				return nil, errors.Join(errors.New("storage config parse fail"), err)
			}
		}
		inventories = append(inventories, inventory)
	}
	return inventories, nil
}

func (s sqliteStorage) GetInventoriesByPolicyAndNamespaceAndHostname(policy, namespace, hostname string) ([]DbInventory, error) {
	selectResult, err := s.db.Query(`
		SELECT id, policy, config, namespace, hostname, name, description, vendor, serial, part_num, type, netbox_id, json_data
		FROM inventories
		WHERE policy = $1 AND namespace = $2 AND hostname = $3
	`, policy, namespace, hostname)
	if err != nil {
		return nil, errors.Join(errors.New("storage fetch inventory fail"), err)
	}
	var inventories []DbInventory
	var configAsString string
	for selectResult.Next() {
		var inventory DbInventory
		err := selectResult.Scan(&inventory.Id, &inventory.Policy, &configAsString, &inventory.Namespace, &inventory.Hostname, &inventory.Name,
			&inventory.Descr, &inventory.Vendor, &inventory.Serial, &inventory.PartNum, &inventory.Type, &inventory.NetboxRefId, &inventory.Blob)
		if err != nil {
			return nil, errors.Join(errors.New("storage create inventory struct fail"), err)
		}
		if len(configAsString) > 0 {
			err = json.Unmarshal([]byte(configAsString), &inventory.Config)
			if err != nil {
				return nil, errors.Join(errors.New("storage config parse fail"), err)
			}
		}
		inventories = append(inventories, inventory)
	}

	return inventories, nil
}

func (s sqliteStorage) UpdateInterface(id string, netboxId int64) (DbInterface, error) {
	_, err := s.db.Exec(`
	UPDATE interfaces SET netbox_id = $1 WHERE id = $2`, netboxId, id)
	if err != nil {
		return DbInterface{}, errors.Join(errors.New("storage update interface fail"), err)
	}
	selectResult := s.db.QueryRow(`
		SELECT id, policy, config, namespace, hostname, name, admin_state, mtu, speed, mac_address, if_type, netbox_id, ip_addresses, json_data
		FROM interfaces
		WHERE id = $1
	`, id)
	var dbInterface DbInterface
	var configAsString string
	var ipsAsString string
	err = selectResult.Scan(&dbInterface.Id, &dbInterface.Policy, &configAsString, &dbInterface.Namespace, &dbInterface.Hostname,
		&dbInterface.Name, &dbInterface.AdminState, &dbInterface.Mtu, &dbInterface.Speed, &dbInterface.MacAddress,
		&dbInterface.IfType, &dbInterface.NetboxRefId, &ipsAsString, &dbInterface.Blob)
	if err != nil {
		return DbInterface{}, errors.Join(errors.New("storage create interface struct fail"), err)
	}
	if len(configAsString) > 0 {
		err = json.Unmarshal([]byte(configAsString), &dbInterface.Config)
		if err != nil {
			return DbInterface{}, errors.Join(errors.New("storage config parse fail"), err)
		}
	}
	err = json.Unmarshal([]byte(ipsAsString), &dbInterface.IpAddresses)
	if err != nil {
		return DbInterface{}, errors.Join(errors.New("storage parse ip address fail"), err)
	}
	return dbInterface, nil
}

func (s sqliteStorage) UpdateDevice(id string, netboxId int64) (DbDevice, error) {
	_, err := s.db.Exec(`
	UPDATE devices SET netbox_id = $1 WHERE id = $2`, netboxId, id)
	if err != nil {
		return DbDevice{}, errors.Join(errors.New("storage update device fail"), err)
	}
	selectResult := s.db.QueryRow(`
		SELECT id, policy, config, namespace, hostname, address, serial_number, model, state, vendor, os, netbox_id, json_data
		FROM devices
		WHERE id = $1`, id)
	var device DbDevice
	var configAsString string
	err = selectResult.Scan(&device.Id, &device.Policy, &configAsString, &device.Namespace, &device.Hostname, &device.Address, &device.SerialNumber,
		&device.Model, &device.State, &device.Vendor, &device.Os, &device.NetboxRefId, &device.Blob)
	if err != nil {
		return DbDevice{}, errors.Join(errors.New("storage create device struct fail"), err)
	}
	if len(configAsString) > 0 {
		err = json.Unmarshal([]byte(configAsString), &device.Config)
		if err != nil {
			return DbDevice{}, errors.Join(errors.New("storage config parse fail"), err)
		}
	}
	return device, nil
}

func (s sqliteStorage) UpdateVlan(id string, netboxId int64) (DbVlan, error) {
	_, err := s.db.Exec(`
	UPDATE vlans SET netbox_id = $1 WHERE id = $2`, netboxId, id)
	if err != nil {
		return DbVlan{}, errors.Join(errors.New("storage update vlan fail"), err)
	}
	selectResult := s.db.QueryRow(`
		SELECT id, policy, config, namespace, hostname, name, state, netbox_id, json_data
		FROM vlans
		WHERE id = $1`, id)
	var vlan DbVlan
	var configAsString string
	err = selectResult.Scan(&vlan.Id, &vlan.Policy, &configAsString, &vlan.Namespace, &vlan.Hostname,
		&vlan.Name, &vlan.State, &vlan.NetboxRefId, &vlan.Blob)
	if err != nil {
		return DbVlan{}, errors.Join(errors.New("storage create vlan struct fail"), err)
	}
	if len(configAsString) > 0 {
		err = json.Unmarshal([]byte(configAsString), &vlan.Config)
		if err != nil {
			return DbVlan{}, errors.Join(errors.New("storage config parse fail"), err)
		}
	}
	return vlan, nil
}

func (s sqliteStorage) UpdateInventory(id string, netboxId int64) (DbInventory, error) {
	_, err := s.db.Exec(`
	UPDATE inventories SET netbox_id = $1 WHERE id = $2`, netboxId, id)
	if err != nil {
		return DbInventory{}, errors.Join(errors.New("storage update inventory fail"), err)
	}
	selectResult := s.db.QueryRow(`
		SELECT id, policy, config, namespace, hostname, name, description, vendor, serial, part_num, type, netbox_id, json_data
		FROM inventories
		WHERE id = $1`, id)
	var inventory DbInventory
	var configAsString string
	err = selectResult.Scan(&inventory.Id, &inventory.Policy, &configAsString, &inventory.Namespace, &inventory.Hostname, &inventory.Name,
		&inventory.Descr, &inventory.Vendor, &inventory.Serial, &inventory.PartNum, &inventory.Type, &inventory.NetboxRefId, &inventory.Blob)
	if err != nil {
		return DbInventory{}, errors.Join(errors.New("storage create inventory struct fail"), err)
	}
	if len(configAsString) > 0 {
		err = json.Unmarshal([]byte(configAsString), &inventory.Config)
		if err != nil {
			return DbInventory{}, errors.Join(errors.New("storage config parse fail"), err)
		}
	}
	return inventory, nil
}

func (s sqliteStorage) Save(policy string, jsonData map[string]interface{}) (stored interface{}, err error) {
	confData := jsonData["config"]
	data, ok := jsonData["interfaces"].([]interface{})
	if ok {
		return s.saveInterfaces(policy, confData, data, err)
	}
	data, ok = jsonData["device"].([]interface{})
	if ok {
		return s.saveDevices(policy, confData, data, err)
	}
	data, ok = jsonData["vlan"].([]interface{})
	if ok {
		return s.saveVlans(policy, confData, data, err)
	}
	data, ok = jsonData["inventory"].([]interface{})
	if ok {
		return s.saveInventories(policy, confData, data, err)
	}
	return nil, errors.New("not able to save anything from entry")
}

func (s sqliteStorage) saveInventories(policy string, conf interface{}, inData []interface{}, err error) (interface{}, error) {
	inventories := make([]DbInventory, len(inData))
	var errs error
	var configAsString string
	if conf != nil {
		b, err := json.Marshal(conf)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			configAsString = string(b)
		}
	}
	for _, inventoryData := range inData {
		dataAsString, err := json.Marshal(inventoryData)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		inventory := DbInventory{
			Id:          uuid.NewString(),
			Config:      conf,
			Policy:      policy,
			NetboxRefId: -1,
			Blob:        string(dataAsString),
		}
		err = json.Unmarshal(dataAsString, &inventory)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		statement, err := s.db.Prepare(
			`INSERT INTO inventories 
					( id, policy, config, namespace, hostname, name, description, vendor, serial, part_num, type, netbox_id, json_data)
				VALUES 
					( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13 )`)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		_, err = statement.Exec(inventory.Id, policy, configAsString, inventory.Namespace, inventory.Hostname, inventory.Name,
			inventory.Descr, inventory.Vendor, inventory.Serial, inventory.PartNum, inventory.Type, inventory.NetboxRefId, dataAsString)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		inventories = append(inventories, inventory)
	}
	return inventories, errs
}

func (s sqliteStorage) saveVlans(policy string, conf interface{}, vData []interface{}, err error) (interface{}, error) {
	vlans := make([]DbVlan, len(vData))
	var errs error
	var configAsString string
	if conf != nil {
		b, err := json.Marshal(conf)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			configAsString = string(b)
		}
	}
	for _, vlanData := range vData {
		dataAsString, err := json.Marshal(vlanData)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		vlan := DbVlan{
			Id:          uuid.NewString(),
			Config:      conf,
			Policy:      policy,
			NetboxRefId: -1,
			Blob:        string(dataAsString),
		}
		err = json.Unmarshal(dataAsString, &vlan)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		statement, err := s.db.Prepare(
			`INSERT INTO vlans 
					( id, policy, config, namespace, hostname, name, state, netbox_id, json_data)
				VALUES 
					( $1, $2, $3, $4, $5, $6, $7, $8, $9 )`)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		_, err = statement.Exec(vlan.Id, policy, configAsString, vlan.Namespace, vlan.Hostname, vlan.Name,
			vlan.State, vlan.NetboxRefId, dataAsString)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		vlans = append(vlans, vlan)
	}
	return vlans, errs
}

func (s sqliteStorage) saveDevices(policy string, conf interface{}, dData []interface{}, err error) (interface{}, error) {
	devicesAdded := make([]DbDevice, len(dData))
	var errs error
	var configAsString string
	if conf != nil {
		b, err := json.Marshal(conf)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			configAsString = string(b)
		}
	}
	for _, deviceData := range dData {
		dataAsString, err := json.Marshal(deviceData)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		dbDevice := DbDevice{
			Id:          uuid.NewString(),
			Policy:      policy,
			Config:      conf,
			NetboxRefId: -1,
			Blob:        string(dataAsString),
		}
		err = json.Unmarshal(dataAsString, &dbDevice)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		statement, err := s.db.Prepare(
			`
				INSERT INTO devices 
					(id, policy, config, namespace, hostname, address, serial_number, model, state, vendor, os, netbox_id, json_data) 
				VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13 )`)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		_, err = statement.Exec(dbDevice.Id, policy, configAsString, dbDevice.Namespace, dbDevice.Hostname, dbDevice.Address, dbDevice.SerialNumber,
			dbDevice.Model, dbDevice.State, dbDevice.Vendor, dbDevice.Os, dbDevice.NetboxRefId, dataAsString)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		devicesAdded = append(devicesAdded, dbDevice)
	}
	return devicesAdded, errs
}

func (s sqliteStorage) saveInterfaces(policy string, conf interface{}, ifData []interface{}, err error) (interface{}, error) {
	interfacesAdded := make([]DbInterface, len(ifData))
	var errs error
	var configAsString string
	if conf != nil {
		b, err := json.Marshal(conf)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			configAsString = string(b)
		}
	}
	for _, interfaceData := range ifData {
		dataAsString, err := json.Marshal(interfaceData)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		dbInterface := DbInterface{
			Id:          uuid.NewString(),
			Policy:      policy,
			Config:      conf,
			NetboxRefId: -1,
			Blob:        string(dataAsString),
		}
		var ipAddresses []IpAddress
		interfaceAsJson := interfaceData.(map[string]interface{})
		ip4AddressList := interfaceAsJson["ipAddressList"].([]interface{})
		for _, ipAddress := range ip4AddressList {
			ipAddresses = append(ipAddresses, IpAddress{
				Address: ipAddress.(string),
				Type:    "v4",
			})
		}
		ip6AddressList := interfaceAsJson["ip6AddressList"].([]interface{})
		for _, ipAddress := range ip6AddressList {
			ipAddresses = append(ipAddresses, IpAddress{
				Address: ipAddress.(string),
				Type:    "v6",
			})
		}
		ipsAsString, err := json.Marshal(ipAddresses)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		dbInterface.IpAddresses = ipAddresses
		err = json.Unmarshal(dataAsString, &dbInterface)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		statement, err := s.db.Prepare(`
			INSERT INTO interfaces 
			    (id, policy, config, namespace, hostname, name, admin_state, mtu, speed, mac_address, if_type, ip_addresses, netbox_id, json_data) 
			VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14 )`)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		_, err = statement.Exec(dbInterface.Id, policy, configAsString, dbInterface.Namespace, dbInterface.Hostname, dbInterface.Name, dbInterface.AdminState,
			dbInterface.Mtu, dbInterface.Speed, dbInterface.MacAddress, dbInterface.IfType, ipsAsString, dbInterface.NetboxRefId, dataAsString)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		interfacesAdded = append(interfacesAdded, dbInterface)
	}
	return interfacesAdded, errs
}

func startSqliteDb(logger *zap.Logger) (db *sql.DB, err error) {
	if !slices.Contains(sql.Drivers(), "sqlite3") {
		logger.Error("SQLite does not have required driver", zap.Error(err))
		return nil, err
	}
	db, err = sql.Open("sqlite3", ":memory")
	if err != nil {
		logger.Error("SQLite could not be initialized", zap.Error(err))
		return nil, err
	}

	createInterfacesTableStatement, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS interfaces 
		( id TEXT PRIMARY KEY, 
		 policy TEXT,
		 config TEXT,
		 namespace TEXT,
		 hostname TEXT,
		 name TEXT,
		 admin_state TEXT,
		 mtu INTEGER,
		 speed INTEGER,
		 mac_address TEXT,
		 if_type TEXT,
		 ip_addresses TEXT,
		 netbox_id INTEGER, 
		 json_data TEXT )`)
	if err != nil {
		logger.Error("error preparing interfaces statement", zap.Error(err))
		return nil, err
	}
	_, err = createInterfacesTableStatement.Exec()
	if err != nil {
		logger.Error("error creating interfaces table", zap.Error(err))
		return nil, err
	}
	logger.Debug("successfully created Interfaces table")
	createDeviceTableStatement, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS devices 
		(
		    id TEXT PRIMARY KEY, 
		 	policy TEXT,
			config TEXT,
		 	namespace TEXT,
		 	hostname TEXT,
		 	serial_number TEXT,
		 	address TEXT,
		 	model TEXT,
		 	state TEXT,
		 	vendor TEXT,
			os TEXT,
		 	netbox_id INTEGER, 
		    json_data TEXT 
		)`)
	if err != nil {
		logger.Error("error preparing devices statement ", zap.Error(err))
		return nil, err
	}
	_, err = createDeviceTableStatement.Exec()
	if err != nil {
		logger.Error("error creating devices table", zap.Error(err))
		return nil, err
	}
	logger.Debug("successfully created devices table")

	createVlansTableStatement, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS vlans
	(
	    id TEXT PRIMARY KEY,
	    policy TEXT,
		config TEXT,
	    namespace TEXT,
		hostname TEXT,
		name TEXT,
		state TEXT,
		netbox_id INTEGER,
		json_data TEXT 
	)`)
	if err != nil {
		logger.Error("error preparing vlans statement ", zap.Error(err))
		return nil, err
	}
	_, err = createVlansTableStatement.Exec()
	if err != nil {
		logger.Error("error creating vlans table", zap.Error(err))
		return nil, err
	}
	createInventoryTableStatement, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS inventories 
		(
		    id TEXT PRIMARY KEY, 
		 	policy TEXT,
			config TEXT,
		 	namespace TEXT,
		 	hostname TEXT,
		 	name TEXT,
		 	description TEXT,
		 	vendor TEXT,
		 	serial TEXT,
		 	part_num TEXT,
			 type TEXT,
		 	netbox_id INTEGER, 
		    json_data TEXT 
		)`)
	if err != nil {
		logger.Error("error preparing inventories statement ", zap.Error(err))
		return nil, err
	}
	_, err = createInventoryTableStatement.Exec()
	if err != nil {
		logger.Error("error creating inventories table", zap.Error(err))
		return nil, err
	}
	logger.Debug("successfully created inventories table")

	constraint1TableStatement, err := db.Prepare(
		`CREATE UNIQUE INDEX IF NOT EXISTS interfaces_uniques ON interfaces(policy, namespace, hostname, name)`)
	if err != nil {
		logger.Error("error constraints statement ", zap.Error(err))
		return nil, err
	}
	_, err = constraint1TableStatement.Exec()
	if err != nil {
		logger.Error("error constraints execution", zap.Error(err))
		return nil, err
	}
	constraint2TableStatement, err := db.Prepare(
		`CREATE UNIQUE INDEX IF NOT EXISTS devices_uniques ON devices(policy, namespace, hostname)`)
	if err != nil {
		logger.Error("error constraints statement ", zap.Error(err))
		return nil, err
	}
	_, err = constraint2TableStatement.Exec()
	if err != nil {
		logger.Error("error constraints execution", zap.Error(err))
		return nil, err
	}
	constraint3TableStatement, err := db.Prepare(
		`CREATE UNIQUE INDEX IF NOT EXISTS vlans_uniques ON vlans(policy, namespace, hostname, name)`)
	if err != nil {
		logger.Error("error constraints statement ", zap.Error(err))
		return nil, err
	}
	_, err = constraint3TableStatement.Exec()
	if err != nil {
		logger.Error("error constraints execution", zap.Error(err))
		return nil, err
	}
	constraint4TableStatement, err := db.Prepare(
		`CREATE UNIQUE INDEX IF NOT EXISTS inventories_uniques ON inventories(policy, namespace, hostname, name)`)
	if err != nil {
		logger.Error("error constraints statement ", zap.Error(err))
		return nil, err
	}
	_, err = constraint4TableStatement.Exec()
	if err != nil {
		logger.Error("error constraints execution", zap.Error(err))
		return nil, err
	}

	return
}
