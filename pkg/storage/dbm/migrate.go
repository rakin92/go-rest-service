package dbm

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

// Migrate runs the migrations scripts in scripts folder
func Migrate(c *cfg.DB) error {
	logger.Info("[Migrate.Scripts] Running DB Migration Scripts")
	mg, err := migrate.New("file://pkg/storage/dbm/scripts/", c.DSN)
	if err != nil {
		return fmt.Errorf("[Migration.Scripts]: %v", err)
	}
	defer mg.Close()

	err = mg.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("[Migration.Scripts]: %v", err)
	}
	logger.Info("[Migrate.Scripts] DB Migration Scripts complete")
	return nil
}
