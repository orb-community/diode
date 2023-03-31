package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

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

func (s sqliteStorage) GetInterfaceByPolicyAndNamespaceAndHostname(policy, namespace, hostname string) ([]DbInterface, error) {
	selectResult, err := s.db.Query(`
		SELECT id, policy, namespace, hostname, name, admin_state, mtu, speed, mac_address, if_type, netbox_id, ip_addresses, json_data
		FROM interfaces
		WHERE policy = $1 AND namespace = $2 AND hostname = $3
	`, policy, namespace, hostname)
	if err != nil {
		return nil, fmt.Errorf("storage - fetch interface fail - %v", err)
	}
	var interfaces []DbInterface
	var ipsAsString string
	for selectResult.Next() {
		var iface DbInterface
		err := selectResult.Scan(&iface.Id, &iface.Policy, &iface.Namespace, &iface.Hostname, &iface.Name, &iface.AdminState,
			&iface.Mtu, &iface.Speed, &iface.MacAddress, &iface.IfType, &iface.NetboxRefId, &ipsAsString, &iface.Blob)
		if err != nil {
			return nil, fmt.Errorf("storage - create ifce struct fail - %v", err)
		}
		err = json.Unmarshal([]byte(ipsAsString), &iface.IpAddresses)
		if err != nil {
			return nil, fmt.Errorf("storage - ip_address parse fail - %v", err)
		}
		interfaces = append(interfaces, iface)
	}

	return interfaces, nil
}

func (s sqliteStorage) GetDevicesByPolicyAndNamespace(policy, namespace string) ([]DbDevice, error) {
	selectResult, err := s.db.Query(`
		SELECT id, policy, namespace, hostname, serial_number, model, state, vendor, netbox_id, json_data
		FROM devices
		WHERE policy = $1 AND namespace = $2
	`, policy, namespace)
	if err != nil {
		return nil, fmt.Errorf("storage - fetch device fail - %v", err)
	}
	var devices []DbDevice
	for selectResult.Next() {
		var device DbDevice
		err := selectResult.Scan(&device.Id, &device.Policy, &device.Namespace, &device.Hostname, &device.SerialNumber,
			&device.Model, &device.State, &device.Vendor, &device.NetboxRefId, &device.Blob)
		if err != nil {
			return nil, fmt.Errorf("storage - create device struct fail - %v", err)
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (s sqliteStorage) GetDeviceByPolicyAndNamespaceAndHostname(policy, namespace, hostname string) (DbDevice, error) {
	selectResult := s.db.QueryRow(`
		SELECT id, policy, namespace, hostname, serial_number, model, state, vendor, netbox_id, json_data
		FROM devices
		WHERE policy = $1 AND namespace = $2 AND hostname = $3
	`, policy, namespace, hostname)
	var device DbDevice
	err := selectResult.Scan(&device.Id, &device.Policy, &device.Namespace, &device.Hostname, &device.SerialNumber,
		&device.Model, &device.State, &device.Vendor, &device.NetboxRefId, &device.Blob)
	if err != nil {
		return DbDevice{}, fmt.Errorf("storage - fetch device fail - %v", err)
	}
	return device, nil
}

func (s sqliteStorage) GetVlansByPolicyAndNamespaceAndHostname(policy, namespace, hostname string) ([]DbVlan, error) {
	selectResult, err := s.db.Query(`
		SELECT id, policy, namespace, hostname, name, state, netbox_id, json_data
		FROM vlans
		WHERE policy = $1 AND namespace = $2 AND hostname = $3
	`, policy, namespace, hostname)
	if err != nil {
		return nil, fmt.Errorf("storage - fetch vlan fail - %v", err)
	}
	var vlans []DbVlan
	for selectResult.Next() {
		var vlan DbVlan
		err := selectResult.Scan(&vlan.Id, &vlan.Policy, &vlan.Namespace, &vlan.Hostname, &vlan.Name,
			&vlan.State, &vlan.NetboxRefId, &vlan.Blob)
		if err != nil {
			return nil, fmt.Errorf("storage - create vlan struct fail - %v", err)
		}
		vlans = append(vlans, vlan)
	}

	return vlans, nil
}

func (s sqliteStorage) UpdateInterface(id string, netboxId int64) (DbInterface, error) {
	_, err := s.db.Exec(`
	UPDATE interfaces SET netbox_id = $1 WHERE id = $2`, netboxId, id)
	if err != nil {
		return DbInterface{}, fmt.Errorf("storage - update interface fail - %v", err)
	}
	selectResult := s.db.QueryRow(`
		SELECT id, policy, namespace, hostname, name, admin_state, mtu, speed, mac_address, if_type, netbox_id, ip_addresses, json_data
		FROM interfaces
		WHERE id = $1
	`, id)
	var dbInterface DbInterface
	var ipsAsString string
	err = selectResult.Scan(&dbInterface.Id, &dbInterface.Policy, &dbInterface.Namespace, &dbInterface.Hostname,
		&dbInterface.Name, &dbInterface.AdminState, &dbInterface.Mtu, &dbInterface.Speed, &dbInterface.MacAddress,
		&dbInterface.IfType, &dbInterface.NetboxRefId, &ipsAsString, &dbInterface.Blob)
	if err != nil {
		return DbInterface{}, fmt.Errorf("storage - create ifce struct fail - %v", err)
	}
	err = json.Unmarshal([]byte(ipsAsString), &dbInterface.IpAddresses)
	if err != nil {
		return DbInterface{}, fmt.Errorf("storage - ip_address parse fail - %v", err)
	}
	return dbInterface, nil
}

func (s sqliteStorage) UpdateDevice(id string, netboxId int64) (DbDevice, error) {
	_, err := s.db.Exec(`
	UPDATE devices SET netbox_id = $1 WHERE id = $2`, netboxId, id)
	if err != nil {
		return DbDevice{}, fmt.Errorf("storage - update device fail - %v", err)
	}
	selectResult := s.db.QueryRow(`
		SELECT id, policy, namespace, hostname, address, serial_number, model, state, vendor, netbox_id, json_data
		FROM devices
		WHERE id = $1`, id)
	var device DbDevice
	err = selectResult.Scan(&device.Id, &device.Policy, &device.Namespace, &device.Hostname, &device.Address, &device.SerialNumber,
		&device.Model, &device.State, &device.Vendor, &device.NetboxRefId, &device.Blob)
	if err != nil {
		return DbDevice{}, fmt.Errorf("storage - create device struct fail - %v", err)
	}
	return device, nil
}

func (s sqliteStorage) UpdateVlan(id string, netboxId int64) (DbVlan, error) {
	_, err := s.db.Exec(`
	UPDATE vlans SET netbox_id = $1 WHERE id = $2`, netboxId, id)
	if err != nil {
		return DbVlan{}, fmt.Errorf("storage - update vlan fail - %v", err)
	}
	selectResult := s.db.QueryRow(`
		SELECT id, policy, namespace, hostname, name, state, netbox_id, json_data
		FROM vlans
		WHERE id = $1`, id)
	var vlan DbVlan
	err = selectResult.Scan(&vlan.Id, &vlan.Policy, &vlan.Namespace, &vlan.Hostname,
		&vlan.Name, &vlan.State, &vlan.NetboxRefId, &vlan.Blob)
	if err != nil {
		return DbVlan{}, fmt.Errorf("storage - create vlan struct fail - %v", err)
	}
	return vlan, nil
}

func (s sqliteStorage) Save(policy string, jsonData map[string]interface{}) (stored interface{}, err error) {
	ifData, ok := jsonData["interfaces"].([]interface{})
	if ok {
		return s.saveInterfaces(policy, ifData, err)
	}
	dData, ok := jsonData["device"].([]interface{})
	if ok {
		return s.saveDevices(policy, dData, err)
	}
	vData, ok := jsonData["vlan"].([]interface{})
	if ok {
		return s.saveVlans(policy, vData, err)
	}
	return nil, errors.New("not able to save anything from entry")
}

func (s sqliteStorage) saveVlans(policy string, vData []interface{}, err error) (interface{}, error) {
	vlans := make([]DbVlan, len(vData))
	for _, vlanData := range vData {
		dataAsString, err := json.Marshal(vlanData)
		if err != nil {
			s.logger.Error("error marshalling interface data", zap.Error(err))
			continue
		}
		vlan := DbVlan{
			Id:          uuid.NewString(),
			Policy:      policy,
			NetboxRefId: -1,
			Blob:        string(dataAsString),
		}
		err = json.Unmarshal(dataAsString, &vlan)
		if err != nil {
			s.logger.Error("error marshalling interface data", zap.Error(err))
			continue
		}
		statement, err := s.db.Prepare(
			`INSERT INTO vlans 
					( id, policy, namespace, hostname, name, state, netbox_id, json_data)
				VALUES 
					( $1, $2, $3, $4, $5, $6, $7, $8 )`)
		if err != nil {
			s.logger.Error("error during preparing insert statement", zap.Error(err))
			continue
		}
		_, err = statement.Exec(vlan.Id, policy, vlan.Namespace, vlan.Hostname, vlan.Name,
			vlan.State, vlan.NetboxRefId, dataAsString)
		if err != nil {
			s.logger.Error("error during preparing insert statement on device",
				zap.Strings("vlan", []string{policy, vlan.Namespace, vlan.Hostname, vlan.Name}),
				zap.Error(err))
			continue
		}
		vlans = append(vlans, vlan)
	}
	return vlans, err
}

func (s sqliteStorage) saveDevices(policy string, dData []interface{}, err error) (interface{}, error) {
	devicesAdded := make([]DbDevice, len(dData))
	for _, deviceData := range dData {
		dataAsString, err := json.Marshal(deviceData)
		if err != nil {
			s.logger.Error("error marshalling interface data", zap.Error(err))
			continue
		}
		dbDevice := DbDevice{
			Id:          uuid.NewString(),
			Policy:      policy,
			NetboxRefId: -1,
			Blob:        string(dataAsString),
		}
		err = json.Unmarshal(dataAsString, &dbDevice)
		if err != nil {
			s.logger.Error("error marshalling interface data", zap.Error(err))
			continue
		}
		statement, err := s.db.Prepare(
			`
				INSERT INTO devices 
					(id, policy, namespace, hostname, address, serial_number, model, state, vendor, netbox_id, json_data) 
				VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11 )`)
		if err != nil {
			s.logger.Error("error during preparing insert statement", zap.Error(err))
			continue
		}
		_, err = statement.Exec(dbDevice.Id, policy, dbDevice.Namespace, dbDevice.Hostname, dbDevice.Address, dbDevice.SerialNumber,
			dbDevice.Model, dbDevice.State, dbDevice.Vendor, dbDevice.NetboxRefId, dataAsString)
		if err != nil {
			s.logger.Error("error during preparing insert statement on device",
				zap.Strings("device", []string{policy, dbDevice.Namespace, dbDevice.Hostname}),
				zap.Error(err))
			continue
		}
		devicesAdded = append(devicesAdded, dbDevice)
	}
	return devicesAdded, err
}

func (s sqliteStorage) saveInterfaces(policy string, ifData []interface{}, err error) (interface{}, error) {
	interfacesAdded := make([]DbInterface, len(ifData))
	for _, interfaceData := range ifData {
		dataAsString, err := json.Marshal(interfaceData)
		if err != nil {
			s.logger.Error("error marshalling interface data", zap.Error(err))
			continue
		}
		dbInterface := DbInterface{
			Id:          uuid.NewString(),
			Policy:      policy,
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
			s.logger.Error("error marshalling ipaddresses from interfaces", zap.Error(err))
			continue
		}
		dbInterface.IpAddresses = ipAddresses
		err = json.Unmarshal(dataAsString, &dbInterface)
		if err != nil {
			s.logger.Error("error marshalling interface data", zap.Error(err))
			continue
		}
		statement, err := s.db.Prepare(`
			INSERT INTO interfaces 
			    (id, policy, namespace, hostname, name, admin_state, mtu, speed, mac_address, if_type, ip_addresses,  netbox_id, json_data) 
			VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13 )`)
		if err != nil {
			s.logger.Error("error during preparing insert statement on interface", zap.Error(err))
			continue
		}
		_, err = statement.Exec(dbInterface.Id, policy, dbInterface.Namespace, dbInterface.Hostname, dbInterface.Name, dbInterface.AdminState,
			dbInterface.Mtu, dbInterface.Speed, dbInterface.MacAddress, dbInterface.IfType, ipsAsString, dbInterface.NetboxRefId, dataAsString)
		if err != nil {
			s.logger.Error("error during preparing insert statement on interface",
				zap.Strings("interface", []string{policy, dbInterface.Namespace, dbInterface.Hostname, dbInterface.Name}),
				zap.Error(err))
			continue
		}
		interfacesAdded = append(interfacesAdded, dbInterface)
	}
	if err != nil {
		return nil, err
	}
	return interfacesAdded, err
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
		 	namespace TEXT,
		 	hostname TEXT,
		 	serial_number TEXT,
		 	address TEXT,
		 	model TEXT,
		 	state TEXT,
		 	vendor TEXT,
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

	return
}
