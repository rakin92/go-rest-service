package migration

import (
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/rakin92/go-rest-service/internal/orm/models"
	"github.com/rakin92/go-rest-service/pkg/consts"
	"github.com/rakin92/go-rest-service/pkg/logger"
)

// updateMigration updates our orm models schemas
func updateMigration(db *gorm.DB) (err error) {
	return db.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.UserProfile{},
		&models.UserAPIKey{},
		&models.User{},
	)
}

// ServiceAutoMigration migrates all the tables and modifications to the connected source
func ServiceAutoMigration(db *gorm.DB) error {
	// Initialize the migration empty so InitSchema runs always first on creation
	m := gormigrate.New(db, gormigrate.DefaultOptions, nil)
	m.InitSchema(func(db *gorm.DB) error {
		logger.Info("[Migration.InitSchema] Initializing database schema")
		switch db.Dialector.Name() {
		case consts.Dialects.PostgresSQL:
			db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
		}
		if err := updateMigration(db); err != nil {
			return fmt.Errorf("[Migration.InitSchema]: %v", err)
		}
		return nil
	})
	m.Migrate()
	if err := updateMigration(db); err != nil {
		return err
	}

	m = gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		SeedRBAC,
		SeedUsers,
	})

	return m.Migrate()
}
