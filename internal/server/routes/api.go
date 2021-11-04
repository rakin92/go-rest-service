package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rakin92/go-rest-service/internal/orm"
	"github.com/rakin92/go-rest-service/internal/server/handlers"
	"github.com/rakin92/go-rest-service/pkg/auth"
	"github.com/rakin92/go-rest-service/pkg/cfg"
)

// AuthAPI is the related routes which is only available user to be authenticated
// user may use weather OAuth with JWT auth token or x-api-key headers
func AuthAPI(sc *cfg.Server, r *gin.Engine, orm *orm.ORM) error {
	// Authorization API group
	authorizedAPI := r.Group(sc.VersionedEndpoint("/api"))
	authorizedAPI.Use(auth.Middleware(sc.VersionedEndpoint("/api"), sc, orm))
	{
		authorizedAPI.GET("/user/:id", handlers.Health())
	}
	return nil
}

// OpenAPI is the related open routes which can be used without being authenticated
func OpenAPI(sc *cfg.Server, r *gin.Engine, orm *orm.ORM) error {
	// Authorization API group
	openAPI := r.Group(sc.VersionedEndpoint("/api"))
	{
		openAPI.GET("/status", handlers.Health())
	}
	return nil
}
