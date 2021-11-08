// Package dbm is a database management package.
// It allows us to interact with sql database via using sqlx.
package dbm

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rakin92/go-rest-service/pkg/cfg"
	"github.com/rakin92/go-rest-service/pkg/logger"
)

type SQLDB struct {
	DB *sqlx.DB
}

// Init creates a connection to our database via sqlx
func Init(c *cfg.DB) (*SQLDB, error) {
	logger.Info("[SQL.Connect] Connecting to %s DB", c.Dialect)
	d, err := sqlx.Connect(c.Dialect, c.DSN)
	if err != nil {
		return nil, err
	}

	if c.MaxIdleCon > 0 {
		d.SetMaxIdleConns(c.MaxIdleCon)
	}
	if c.MaxCon > 0 {
		d.SetMaxOpenConns(c.MaxCon)
	}
	d.SetConnMaxLifetime(time.Hour)
	logger.Info("[SQL.Connect] Connected to %s DB", c.Dialect)
	return &SQLDB{DB: d}, nil
}
