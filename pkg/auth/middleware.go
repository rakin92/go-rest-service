package auth

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/rakin92/travel/internal/orm"
	"github.com/rakin92/travel/pkg/cfg"
	"github.com/rakin92/travel/pkg/logger"

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
	logger.Infof("[Auth.Middleware] Applied to path: %s", path)
	return gin.HandlerFunc(func(c *gin.Context) {
		if a, err := ParseAPIKey(c, cfg); err == nil {
			user, err := orm.FindUserByAPIKey(a)
			if err != nil {
				authError(c, ErrForbidden)
			}
			if user != nil {
				c.Request = addUserIdToContext(c, user.ID)
				logger.Debugf("User authenticated via api: %s", user.ID)
			}
			c.Next()
		} else {
			if err != ErrEmptyAPIKeyHeader {
				authError(c, err)
			} else {
				t, err := ParseToken(c, cfg)
				if err != nil {
					authError(c, err)
				} else {
					if claims, ok := t.Claims.(jwt.MapClaims); ok {
						if claims["exp"] != nil {
							issuer := claims["iss"].(string)
							userid := claims["jti"].(string)
							email := claims["email"].(string)
							if claims["aud"] != nil {
								audiences := claims["aud"]
								logger.Warnf("\n\naudiences: %s\n\n", audiences)
							}
							if claims["alg"] != nil {
								algo := claims["alg"].(string)
								logger.Warnf("\n\nalgo: %s\n\n", algo)
							}
							if user, err := orm.FindUserByJWT(email, issuer, userid); err != nil {
								authError(c, ErrForbidden)
							} else {
								if user != nil {
									c.Request = addUserIdToContext(c, user.ID)
									logger.Debugf("User: %s", user.ID)
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
