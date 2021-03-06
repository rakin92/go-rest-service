// Package orm provides `GORM` helpers for the creation, migration and access
// on the project's database
package orm

import (
	"errors"
	"fmt"

	"github.com/markbates/goth"
	"github.com/rakin92/go-rest-service/internal/orm/migration"
	"github.com/rakin92/go-rest-service/internal/orm/models"
	"github.com/rakin92/go-rest-service/pkg/cfg"
	"github.com/rakin92/go-rest-service/pkg/consts"
	"github.com/rakin92/go-rest-service/pkg/logger"
	"github.com/rakin92/go-rest-service/pkg/storage/dbm"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	sUserTbl  = "User"
	nestedFmt = "%s.%s"
)

// ORM struct to holds the gorm pointer to db
type ORM struct {
	DB *gorm.DB
}

// Init creates a db connection with the selected dialect and connection
// along with running all the db migrations
func Init(c *cfg.DB) (*ORM, error) {
	db, err := gorm.Open(postgres.Open(c.DSN))
	if err != nil {
		logger.Panic(&err, "[ORM] err: %s", err.Error())
	}
	orm := &ORM{DB: db}
	// Log every SQL command on dev, @prod: this should be disabled? Maybe.
	// db.LogMode(c.LogMode) TODO: look into this
	// Automigrate tables
	if c.AutoMigrate {
		// migrates our sql scripts
		err = dbm.Migrate(c)
		if err != nil {
			logger.Fatal(&err, "[sql.Migrate] scripts err: %s", err.Error())
		}
		// migrats our orm schema
		err = migration.ServiceAutoMigration(orm.DB)
		if err != nil {
			logger.Fatal(&err, "[ORM.autoMigrate] err: %v", err.Error())
		}
	}

	logger.Info("[ORM] Database connection initialized.")
	return orm, nil
}

//FindUserByAPIKey finds the user that is related to the API key
func (o *ORM) FindUserByAPIKey(apiKey string) (*models.User, error) {
	if apiKey == "" {
		return nil, errors.New("API key is empty")
	}
	uak := &models.UserAPIKey{}
	usrPerm := fmt.Sprintf(nestedFmt, sUserTbl, consts.EntityNames.Permissions)
	usrRole := fmt.Sprintf(nestedFmt, sUserTbl, consts.EntityNames.Roles)
	if err := o.DB.Preload(sUserTbl).Preload(usrPerm).Preload(usrRole).
		First(uak, "api_key = ?", apiKey).Error; err != nil {
		return nil, err
	}
	return &uak.User, nil
}

// FindUserByJWT finds the user that is related to the APIKey token
func (o *ORM) FindUserByJWT(email string, provider string, userID string) (*models.User, error) {
	if provider == "" || userID == "" {
		return nil, errors.New("provider or userId empty")
	}
	tx := o.DB.Begin()
	p := &models.UserProfile{}
	usrPerm := fmt.Sprintf(nestedFmt, sUserTbl, consts.EntityNames.Permissions)
	usrRole := fmt.Sprintf(nestedFmt, sUserTbl, consts.EntityNames.Roles)
	if err := tx.Preload(sUserTbl).Preload(usrPerm).Preload(usrRole).
		First(p, "email  = ? AND provider = ? AND external_user_id = ?", email, provider, userID).Error; err != nil {
		return nil, err
	}
	return &p.User, nil
}

// UpsertUserProfile saves the user if doesn't exists and adds the OAuth profile
// and updates existing user info if record exist in db
func (o *ORM) UpsertUserProfile(gu *goth.User) (*models.User, error) {
	db := o.DB
	up := &models.UserProfile{}
	u, err := models.GothUserToDBUser(gu, false)
	if err != nil {
		return nil, err
	}

	tx := db.Where("email = ?", gu.Email).First(u)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}

	if tx := db.Model(u).Save(u); tx.Error != nil {
		return nil, tx.Error
	}

	tx = db.Where("email = ? AND provider = ? AND external_user_id = ?", gu.Email, gu.Provider, gu.UserID).First(up)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}

	up, err = models.GothUserToDBUserProfile(gu, false)
	if err != nil {
		return nil, err
	}

	up.User = *u
	if tx := db.Model(up).Save(up); tx.Error != nil {
		return nil, tx.Error
	}
	return u, nil
}
