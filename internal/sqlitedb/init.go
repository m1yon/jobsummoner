package sqlitedb

import (
	"database/sql"

	"github.com/pkg/errors"
)

const (
	ErrOpeningDB = "problem opening db"
	ErrPingingDB = "db did not respond to ping"
)

func NewDB(driverName string, dataSourceName string, open func(driverName string, dataSourceName string) (*sql.DB, error)) (*sql.DB, error) {
	db, err := open(driverName, dataSourceName)

	if err != nil {
		return nil, errors.Wrap(err, ErrOpeningDB)
	}

	return db, nil
}
