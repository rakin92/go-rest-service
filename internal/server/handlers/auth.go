package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/markbates/goth/gothic"

	"github.com/rakin92/go-rest-service/internal/orm"
	"github.com/rakin92/go-rest-service/pkg/cfg"
	"github.com/rakin92/go-rest-service/pkg/consts"
	"github.com/rakin92/go-rest-service/pkg/logger"
)

// AuthProviders begin login with the auth provider
func AuthProviders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// You have to add value context with provider name to get provider name in GetProviderName method
		c.Request = addProviderToContext(c, c.Param(string(consts.ProjectContextKeys.ProviderCtxKey)))
		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err != nil {
			gothic.BeginAuthHandler(c.Writer, c.Request)
		} else {
			logger.Debug("user: %#v", gothUser)
		}
	}
}

// Callback callback to complete auth provider flow
func Callback(sc *cfg.Server, orm *orm.ORM) gin.HandlerFunc {
	return func(c *gin.Context) {
		// You have to add value context with provider name to get provider name in GetProviderName method
		c.Request = addProviderToContext(c, c.Param(string(consts.ProjectContextKeys.ProviderCtxKey)))
		gothUsr, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// finds the user in our system from goth jwt
		u, err := orm.FindUserByJWT(gothUsr.Email, gothUsr.Provider, gothUsr.UserID)
		if err != nil {
			if u, err = orm.UpsertUserProfile(&gothUsr); err != nil {
				logger.Error(&err, "[Auth.CallBack.UserLoggedIn.FindUserByJWT.Error]: %s", err.Error())
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
		claims := &jwt.RegisteredClaims{
			ID:        gothUsr.UserID,
			Subject:   u.Email,
			Issuer:    gothUsr.Provider,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(gothUsr.ExpiresAt.UTC()),
		}

		// issue a new JWT token
		jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod(sc.JWT.Algorithm), claims)
		token, err := jwtToken.SignedString([]byte(sc.JWT.Secret))
		if err != nil {
			logger.Error(&err, "[Auth.Callback.JWT] error: %s", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		json := gin.H{
			"type":          "Bearer",
			"token":         token,
			"refresh_token": gothUsr.RefreshToken,
			"user_id":       u.ID,
		}
		c.JSON(http.StatusOK, json)
	}
}

// Logout logs out of the auth provider
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = addProviderToContext(c, c.Param(string(consts.ProjectContextKeys.ProviderCtxKey)))
		gothic.Logout(c.Writer, c.Request)
		c.Writer.Header().Set("Location", "/")
		c.Writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}

// addProviderToContext adds our auth providers to context
func addProviderToContext(c *gin.Context, value interface{}) *http.Request {
	c.Set(string(consts.ProjectContextKeys.GothicProviderCtxKey), value)
	return c.Request.WithContext(context.WithValue(c.Request.Context(),
		consts.ProjectContextKeys.GothicProviderCtxKey, value))
}
