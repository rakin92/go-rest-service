package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rakin92/go-rest-service/internal/orm"
	"github.com/rakin92/go-rest-service/internal/server/handlers"
	"github.com/rakin92/go-rest-service/pkg/cfg"
	"github.com/rakin92/go-rest-service/pkg/consts"
)

// Auth routes
func Auth(sc *cfg.Server, r *gin.Engine, orm *orm.ORM) error {
	provider := string(consts.ProjectContextKeys.ProviderCtxKey)
	// OAuth handlers
	rg := r.Group(sc.VersionedEndpoint("/auth"))
	rg.GET("/:"+provider, handlers.AuthProviders())
	rg.GET("/:"+provider+"/callback", handlers.Callback(sc, orm))

	return nil
}
