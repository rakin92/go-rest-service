package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rakin92/go-rest-service/internal/orm"
	"github.com/rakin92/go-rest-service/internal/server/handlers"
	"github.com/rakin92/go-rest-service/pkg/auth"
	"github.com/rakin92/go-rest-service/pkg/cfg"
)

// Misc routes
func Misc(sc *cfg.Server, r *gin.Engine, orm *orm.ORM) error {
	// Simple keep-alive/ping handler
	r.GET(sc.VersionedEndpoint("/health"), handlers.Health())
	r.GET(sc.VersionedEndpoint("/secure-health"),
		auth.Middleware(sc.VersionedEndpoint("/secure-health"), sc, orm), handlers.Health())
	return nil
}
