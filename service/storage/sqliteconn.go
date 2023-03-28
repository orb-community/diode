package storage

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

func (s sqliteStorage) Save(policy string, jsonData map[string]interface{}) (id string, err error) {
	data, ok := jsonData["interfaces"].([]map[string]interface{})
	if ok {
		for _, interfaceData := range data {
			dataAsString, err := json.Marshal(interfaceData)
			if err != nil {
				s.logger.Error("error marshalling interface data", zap.Error(err))
				continue
			}
			id := uuid.NewString()
			statement, err := s.db.Prepare("INSERT INTO interfaces (id, policy, interfaceData) VALUES (?,?,?)")
			if err != nil {
				s.logger.Error("error during preparing insert statement on interface", zap.Error(err))
				continue
			}
			statement.Exec(id, policy, dataAsString)
		}
		if err != nil {
			return "", err
		}
	}
	data, ok = jsonData["device"].([]map[string]interface{})
	dataAsString, err := json.Marshal(data)
	if err != nil {
		s.logger.Error("error marshalling interface data", zap.Error(err))
		return "", err
	}
	id = uuid.NewString()
	statement, err := s.db.Prepare("INSERT INTO devices (id, policy, interfaceData) VALUES (?,?,?)")
	if err != nil {
		s.logger.Error("error during preparing insert statement", zap.Error(err))
		return "", err
	}
	_, err = statement.Exec(id, policy, dataAsString)
	if err != nil {
		s.logger.Error("error during executing insert statement", zap.Error(err))
		return "", err
	}

	return id, nil
}

func startSqliteDb(logger *zap.Logger) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", ":memory")
	if err != nil {
		logger.Error("SQLite could not be initialized", zap.Error(err))
		return nil, err
	}

	createInterfacesTableStatement, err := db.Prepare("CREATE TABLE IF NOT EXISTS interfaces (id TEXT PRIMARY KEY, policy TEXT, interfaceData TEXT )")
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
	createDeviceTableStatement, err := db.Prepare("CREATE TABLE IF NOT EXISTS devices (id TEXT PRIMARY KEY, policy TEXT, deviceData TEXT )")
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

	return
}
