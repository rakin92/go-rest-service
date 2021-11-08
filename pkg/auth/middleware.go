package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rakin92/go-rest-service/internal/orm"
	"github.com/rakin92/go-rest-service/pkg/cfg"
	"github.com/rakin92/go-rest-service/pkg/consts"
	"github.com/rakin92/go-rest-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

func authError(c *gin.Context, err error) {
	errKey := "message"
	errMsgHeader := "[Auth] error: "
	e := gin.H{errKey: errMsgHeader + err.Error()}
	c.AbortWithStatusJSON(http.StatusUnauthorized, e)
}

// Middleware wraps the request with auth middleware
func Middleware(path string, cfg *cfg.Server, orm *orm.ORM) gin.HandlerFunc {
	logger.Info("[Auth.Middleware] Applied to path: %s", path)
	return gin.HandlerFunc(func(c *gin.Context) {
		// Check and authenticate with api key
		if a, err := ParseAPIKey(c, cfg); err == nil {
			user, err := orm.FindUserByAPIKey(a)
			if err != nil {
				authError(c, ErrForbidden)
			}
			if user != nil {
				c.Request = addToContext(c, consts.ProjectContextKeys.UserCtxKey, user)
				c.Request = addUserIdToContext(c, user.ID)
				logger.Debug("User authenticated via api: %s", user.ID)
			}
			c.Next()
		} else {
			if err != ErrEmptyAPIKeyHeader {
				authError(c, err)
			} else {
				// Authenticate via JWT Token
				t, err := ParseToken(c, cfg)
				if err != nil {
					authError(c, err)
				} else {
					if claims, ok := t.Claims.(jwt.MapClaims); ok {
						logger.Info("%+v", claims)
						if claims["exp"] != nil {
							issuer := claims["iss"].(string)
							userid := claims["jti"].(string)
							email := claims["sub"].(string)
							if claims["aud"] != nil {
								audiences := claims["aud"]
								logger.Warn("\n\naudiences: %s\n\n", audiences)
							}
							if claims["alg"] != nil {
								algo := claims["alg"].(string)
								logger.Warn("\n\nalgo: %s\n\n", algo)
							}
							if user, err := orm.FindUserByJWT(email, issuer, userid); err != nil {
								authError(c, ErrForbidden)
							} else {
								if user != nil {
									c.Request = addToContext(c, consts.ProjectContextKeys.UserCtxKey, user)
									c.Request = addUserIdToContext(c, user.ID)
									logger.Debug("User: %s", user.ID)
								}
								c.Next()
							}
						} else {
							authError(c, ErrMissingExpField)
						}
					} else {
						authError(c, err)
					}
				}
			}
		}
	})
}
