package database

import "database/sql"

type ConnectionOpener interface {
	Open(driverName string, dataSourceName string) (*sql.DB, error)
}

type SqlConnectionOpener struct{}

func (s *SqlConnectionOpener) Open(driverName string, dataSourceName string) (*sql.DB, error) {
	return sql.Open(driverName, dataSourceName)
}
