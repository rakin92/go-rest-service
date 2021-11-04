// Service command main package lets us start out service
package main

import (
	"github.com/rakin92/go-rest-service/internal/orm"
	"github.com/rakin92/go-rest-service/internal/server"
	"github.com/rakin92/go-rest-service/pkg/cfg"
	"github.com/rakin92/go-rest-service/pkg/env"
	"github.com/rakin92/go-rest-service/pkg/logger"
)

// main function starts our service by initializing our server configuration
// running on db migrations and establashing connection to db
// starts our server with gin
func main() {
	// Initializes our server config loading environment variables
	var conf = &cfg.Server{
		ServiceName:    env.MustGet("SERVICE_NAME"),
		Version:        env.MustGet("APP_VERSION"),
		Env:            env.MustGet("APP_ENV"),
		Host:           env.MustGet("SERVER_HOST"),
		Port:           env.MustGet("SERVER_PORT"),
		URISchema:      env.MustGet("SERVER_URI_SCHEMA"),
		ServiceVersion: env.MustGet("SERVER_PATH_VERSION"),
		SessionSecret:  env.MustGet("SESSION_SECRET"),
		JWT: cfg.JWT{
			Secret:    env.MustGet("AUTH_JWT_SECRET"),
			Algorithm: env.MustGet("AUTH_JWT_SIGNING_ALGORITHM"),
		},
		Database: cfg.DB{
			Dialect:     env.MustGet("GORM_DIALECT"),
			DSN:         env.MustGet("GORM_CONNECTION_DSN"),
			SeedDB:      env.MustGetBool("GORM_SEED_DB"),
			LogMode:     env.MustGetBool("GORM_LOGMODE"),
			AutoMigrate: env.MustGetBool("GORM_AUTOMIGRATE"),
		},
		AuthProviders: []cfg.AuthProvider{
			{
				Provider:  "facebook",
				ClientKey: env.MustGet("PROVIDER_FACEBOOK_KEY"),
				Secret:    env.MustGet("PROVIDER_FACEBOOK_SECRET"),
			},
			{
				Provider:  "google",
				ClientKey: env.MustGet("PROVIDER_GOOGLE_KEY"),
				Secret:    env.MustGet("PROVIDER_GOOGLE_SECRET"),
			},
			{
				Provider:  "twitter",
				ClientKey: env.MustGet("PROVIDER_TWITTER_KEY"),
				Secret:    env.MustGet("PROVIDER_TWITTER_SECRET"),
			},
		},
	}

	// Initialize our database orm
	o, err := orm.Factory(&conf.Database)
	if err != nil {
		logger.Panic(&err, "[ORM]: Failed to connect to database: %s", err.Error())
	}

	// Runs the gin service
	server.Run(conf, o)
}
