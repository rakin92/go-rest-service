package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rakin92/go-rest-service/internal/orm"
	"github.com/rakin92/go-rest-service/internal/server/routes"
	"github.com/rakin92/go-rest-service/pkg/cfg"
	"github.com/rakin92/go-rest-service/pkg/logger"
	"github.com/rakin92/go-rest-service/pkg/storage/cache"
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

	// Authenticated API routes
	if err = routes.AuthAPI(sc, r, orm); err != nil {
		return err
	}

	// Open API routes
	if err = routes.OpenAPI(sc, r, orm); err != nil {
		return err
	}

	return err
}

// Run spins up the server
func Run(sc *cfg.Server, orm *orm.ORM, che *cache.Cache) {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(logger.Middleware(sc.ServiceName))

	// Initialize the Auth providers
	initializeAuthProviders(sc)

	// Routes and Handlers
	registerRoutes(sc, r, orm)

	// Inform the user where the server is listening
	logger.Info("Running %s @ %s", sc.ServiceName, sc.SchemaVersionedEndpoint(""))

	// Run the server
	// Print out and exit(1) to the OS if the server cannot run
	if err := r.Run(sc.ListenEndpoint()); err != nil {
		logger.Fatal(&err, "Failed to start service %s", sc.ServiceName)
	}
}
