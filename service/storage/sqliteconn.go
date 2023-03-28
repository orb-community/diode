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
	next   Service
}

func NewSqliteStorage(db *sql.DB, logger *zap.Logger) Service {
	return sqliteStorage{db: db, logger: logger}
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
				s.logger.Error("error during preparing insert statementon interface", zap.Error(err))
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
	statement.Exec(id, policy, dataAsString)

	return s.next.Save(policy, jsonData)
}
