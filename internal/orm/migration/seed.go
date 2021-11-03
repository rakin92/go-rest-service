package migration

import (
	"reflect"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gofrs/uuid"
	"github.com/rakin92/travel/internal/orm/models"
	"github.com/rakin92/travel/pkg/consts"
	"github.com/rakin92/travel/pkg/logger"
	"gorm.io/gorm"
)

func rollback(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
	}
}

var (
	fname = "Test"
	lname = "User"
	users = []*models.User{
		{
			Email:     "admin@test.com",
			FirstName: &fname,
			LastName:  &lname,
			Roles:     []models.Role{{BaseModelSeq: models.BaseModelSeq{ID: 1}}},
		},
		{
			Email:     "user@test.com",
			FirstName: &fname,
			LastName:  &lname,
			Roles:     []models.Role{{BaseModelSeq: models.BaseModelSeq{ID: 2}}},
		},
	}
)

// SeedUsers inserts the first users
var SeedUsers *gormigrate.Migration = &gormigrate.Migration{
	ID: "SEED_USERS",
	Migrate: func(db *gorm.DB) error {
		tx := db.Begin()
		defer rollback(tx)
		for _, u := range users {
			if err := tx.Create(u).Error; err != nil {
				return err
			}
			apiKey := uuid.Must(uuid.NewV4()).String()
			if err := tx.Create(&models.UserAPIKey{UserID: u.ID, APIKey: apiKey}).Error; err != nil {
				return err
			}
		}
		tx.Commit()
		return nil
	},
	Rollback: func(db *gorm.DB) error {
		tx := db.Begin()
		defer rollback(tx)
		for _, u := range users {
			if err := tx.Delete(u).Error; err != nil {
				return err
			}
		}
		tx.Commit()
		return nil
	},
}

// SeedRBAC inserts the first users
var SeedRBAC *gormigrate.Migration = &gormigrate.Migration{
	ID: "SEED_RBAC",
	Migrate: func(db *gorm.DB) error {
		tx := db.Begin()
		defer rollback(tx)
		v := reflect.ValueOf(consts.EntityNames)
		tablenames := make([]interface{}, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			tablenames[i] = consts.GetTableName(v.Field(i).Interface().(string))
		}
		v = reflect.ValueOf(consts.Permissions)
		permissions := make([]interface{}, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			permissions[i] = v.Field(i).Interface()
		}
		padmin := []models.Permission{}
		puser := []models.Permission{}
		for _, t := range tablenames {
			for _, p := range permissions {
				permission := models.Permission{
					Tag:         consts.FormatPermissionTag(p.(string), t.(string)),
					Description: consts.FormatPermissionDesc(p.(string), t.(string)),
				}
				if err := tx.Create(&permission).First(&permission).Error; err != nil {
					logger.Errorf(&err, "[Migration.Jobs.SeedRBAC.permissions] error: %s", err.Error())
					return err
				}
				padmin = append(padmin, permission)
				puser = append(puser, permission)
			}
		}
		for _, r := range consts.Roles {
			role := &models.Role{
				Name:        r.Name,
				Description: r.Description,
			}
			if err := tx.Create(role).First(&role).Error; err != nil {
				logger.Errorf(&err, "[Migration.Jobs.SeedRBAC.roles] error: %s", err.Error())
				return err
			}
			switch r.Name {
			case "admin":
				for _, p := range padmin {
					tx.Model(role).Association(consts.EntityNames.Permissions).Append(p)
				}
			case "user":
				for _, p := range puser {
					tx.Model(role).Association(consts.EntityNames.Permissions).Append(p)
				}
			}
		}
		tx.Commit()
		return nil
	},
	Rollback: func(db *gorm.DB) error {
		tx := db.Begin()
		defer rollback(tx)
		for _, u := range users {
			if err := tx.Delete(u).Error; err != nil {
				return err
			}
		}
		tx.Commit()
		return nil
	},
}
