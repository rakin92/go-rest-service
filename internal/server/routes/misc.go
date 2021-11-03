package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rakin92/travel/internal/orm"
	"github.com/rakin92/travel/internal/server/handlers"
	"github.com/rakin92/travel/pkg/auth"
	"github.com/rakin92/travel/pkg/cfg"
)

// Misc routes
func Misc(sc *cfg.Server, r *gin.Engine, orm *orm.ORM) error {
	// Simple keep-alive/ping handler
	r.GET(sc.VersionedEndpoint("/health"), handlers.Health())
	r.GET(sc.VersionedEndpoint("/secure-health"),
		auth.Middleware(sc.VersionedEndpoint("/secure-health"), sc, orm), handlers.Health())
	return nil
}
