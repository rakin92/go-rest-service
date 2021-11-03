package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rakin92/travel/internal/orm"
	"github.com/rakin92/travel/pkg/auth"
	"github.com/rakin92/travel/pkg/cfg"
)

func API(sc *cfg.Server, r *gin.Engine, orm *orm.ORM) error {
	// Authorization API group
	authorizedAPI := r.Group(sc.VersionedEndpoint("/api"))
	authorizedAPI.Use(auth.Middleware(sc.VersionedEndpoint("/api"), sc, orm))
	return nil
}
