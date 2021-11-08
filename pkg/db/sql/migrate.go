package sql

import (
	"fmt"

	// postgres required for golang-migrate db dialact
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// file required for golang-migrate file system
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
	"github.com/rakin92/go-rest-service/pkg/cfg"
	"github.com/rakin92/go-rest-service/pkg/logger"
)

// MigrateScripts runs the migrations scripts in scripts folder
func Migrate(c *cfg.DB) error {
	logger.Info("[Migration.Scripts] Running DB Migration Scripts")
	mg, err := migrate.New("file://pkg/db/sql/scripts/", c.DSN)
	if err != nil {
		return fmt.Errorf("[Migration.Scripts]: %v", err)
	}
	defer mg.Close()

	err = mg.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("[Migration.Scripts]: %v", err)
	}
	return nil
}
