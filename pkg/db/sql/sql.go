package sql

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rakin92/go-rest-service/pkg/cfg"
)

type SQLDB struct {
	DB *sqlx.DB
}

func Init(c *cfg.DB) (*SQLDB, error) {
	db, err := sqlx.Connect(c.Dialect, c.DSN)
	if err != nil {
		return nil, err
	}

	if c.MaxIdleCon > 0 {
		db.SetMaxIdleConns(c.MaxIdleCon)
	}
	if c.MaxCon > 0 {
		db.SetMaxOpenConns(c.MaxCon)
	}
	db.SetConnMaxLifetime(time.Hour)

	return &SQLDB{DB: db}, nil
}
