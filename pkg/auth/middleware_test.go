// Package auth holds some standard auth functionalities.
// It provides us middleware easily into our routers,
// and manage access control of our apis.
package auth

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rakin92/go-rest-service/internal/orm"
	"github.com/rakin92/go-rest-service/internal/orm/models"
	"github.com/rakin92/go-rest-service/pkg/cfg"
)

func TestMiddleware(t *testing.T) {
	// TODO Need to update test
	t.Run("middleware api key from header", func(t *testing.T) {
		findUserByAPIKey = func(apiKey string, o *orm.ORM) (*models.User, error) {
			return &models.User{}, nil
		}
		parseAPIKey = func(c *gin.Context, sc *cfg.Server) (apiKey string, err error) {
			return "api_key", nil
		}
		svc := &cfg.Server{}
		Middleware("/test", svc, &orm.ORM{})
	})
	// TODO Need to update test
	t.Run("middleware jwt key from header", func(t *testing.T) {
		findUserByJWT = func(email, provider, userID string, o *orm.ORM) (*models.User, error) {
			return &models.User{}, nil
		}

		parseToken = func(c *gin.Context, sc *cfg.Server) (t *jwt.Token, err error) {
			return &jwt.Token{}, nil
		}
		svc := &cfg.Server{}
		Middleware("/test", svc, &orm.ORM{})
	})
}
