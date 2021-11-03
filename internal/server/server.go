package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rakin92/travel/internal/orm"
	"github.com/rakin92/travel/internal/server/routes"
	"github.com/rakin92/travel/pkg/cfg"
	"github.com/rakin92/travel/pkg/logger"
)

// registerRoutes register the routes for the server
func registerRoutes(sc *cfg.Server, r *gin.Engine, orm *orm.ORM) (err error) {

	// Miscellaneous routes
	if err = routes.Misc(sc, r, orm); err != nil {
		return err
	}

	// Auth routes
	if err = routes.Auth(sc, r, orm); err != nil {
		return err
	}

	// API routes
	if err = routes.API(sc, r, orm); err != nil {
		return err
	}

	return err
}

// Run spins up the server
func Run(sc *cfg.Server, orm *orm.ORM) {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(logger.Middleware(sc.ServiceName))

	// Initialize the Auth providers
	initalizeAuthProviders(sc)

	// Routes and Handlers
	registerRoutes(sc, r, orm)

	// Inform the user where the server is listening
	logger.Infof("Running %s @ %s", sc.ServiceName, sc.SchemaVersionedEndpoint(""))

	// Run the server
	// Print out and exit(1) to the OS if the server cannot run
	if err := r.Run(sc.ListenEndpoint()); err != nil {
		logger.Fatalf(&err, "Failed to start service %s", sc.ServiceName)
	}
}
