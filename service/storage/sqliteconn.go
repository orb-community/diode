package storage

import "database/sql"

type sqliteStorage struct {
	db *sql.DB
}

func NewSqliteStorage(db *sql.DB) Service {
	return sqliteStorage{db: db}
}

func (s sqliteStorage) StoreLog() {

}
