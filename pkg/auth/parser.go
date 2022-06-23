package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"

	"github.com/rakin92/go-rest-service/pkg/cfg"
	"github.com/rakin92/go-rest-service/pkg/consts"
)

var (
	// APIKeyHeader The API key header name
	APIKeyHeader = "x-api-key"

	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName = "Bearer"

	// APIKeyLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	APIKeyLookup = "param:api_key,query:api_key,cookie:api_key,header:" + APIKeyHeader

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup = "param:api_key,query:token,cookie:jwt,header:Authorization"

	// ErrNoClaims when HTTP status 403 is given
	ErrNoClaims = errors.New("invalid token")

	// ErrForbidden when HTTP status 403 is given
	ErrForbidden = errors.New("you don't have permission to access this resource")

	// ErrExpiredToken indicates JWT token has expired. Can't refresh.
	ErrExpiredToken = errors.New("token is expired")

	// ErrEmptyAuthHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
	ErrEmptyAuthHeader = errors.New("auth header is empty")

	// ErrEmptyAPIKeyHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
	ErrEmptyAPIKeyHeader = errors.New("api key header is empty")

	// ErrMissingExpField missing exp field in token
	ErrMissingExpField = errors.New("missing exp field")

	// ErrInvalidAuthHeader indicates auth header is invalid, could for example have the wrong Realm name
	ErrInvalidAuthHeader = errors.New("auth header is invalid")

	// ErrEmptyQueryToken can be thrown if authing with URL Query, the query token variable is empty
	ErrEmptyQueryToken = errors.New("query token is empty")

	// ErrEmptyCookieToken can be thrown if authing with a cookie, the token cokie is empty
	ErrEmptyCookieToken = errors.New("cookie token is empty")

	// ErrEmptyParamToken can be thrown if authing with parameter in path, the parameter in path is empty
	ErrEmptyParamToken = errors.New("parameter token is empty")

	// ErrInvalidSigningAlgorithm indicates signing algorithm is invalid, needs to be HS256, HS384, HS512, RS256, RS384 or RS512
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")

	jwtParse            = jwt.Parse
	jwtGetSigningMethod = jwt.GetSigningMethod
)

// jwtFromHeader retrieves jwt token from header
func jwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		return "", ErrEmptyAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == TokenHeadName) {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}

func apiKeyFromHeader(c *gin.Context, key string) (string, error) {
	apiKey := c.Request.Header.Get(key)
	if apiKey == "" {
		return "", ErrEmptyAPIKeyHeader
	}
	return apiKey, nil
}

func tokenFromQuery(c *gin.Context, key string) (string, error) {
	token := c.Query(key)
	if token == "" {
		return "", ErrEmptyQueryToken
	}
	return token, nil
}

func tokenFromCookie(c *gin.Context, key string) (string, error) {
	cookie, _ := c.Cookie(key)
	if cookie == "" {
		return "", ErrEmptyCookieToken
	}
	return cookie, nil
}

func tokenFromParam(c *gin.Context, key string) (string, error) {
	token := c.Param(key)
	if token == "" {
		return "", ErrEmptyParamToken
	}
	return token, nil
}

// ParseToken parse jwt token from gin context
// looks for token in header, query params, cookie
func ParseToken(c *gin.Context, sc *cfg.Server) (t *jwt.Token, err error) {
	var token string
	methods := strings.Split(TokenLookup, ",")
	for _, method := range methods {
		if len(token) > 0 {
			break
		}
		parts := strings.Split(strings.TrimSpace(method), ":")
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		switch k {
		case "header":
			token, err = jwtFromHeader(c, v)
		case "query":
			token, err = tokenFromQuery(c, v)
		case "cookie":
			token, err = tokenFromCookie(c, v)
		case "param":
			token, err = tokenFromParam(c, v)
		}
	}
	if err != nil {
		return nil, err
	}
	SigningAlgorithm := sc.JWT.Algorithm
	Key := []byte(sc.JWT.Secret)
	return jwtParse(token, func(t *jwt.Token) (any, error) {
		if jwtGetSigningMethod(SigningAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}
		c.Set("AUTH_JWT_TOKEN", token)
		return Key, nil
	})
}

// ParseAPIKey parse api key from gin context
// looks for x-api-key in header, query params, cookie
func ParseAPIKey(c *gin.Context, sc *cfg.Server) (apiKey string, err error) {
	methods := strings.Split(APIKeyLookup, ",")
	for _, method := range methods {
		if len(apiKey) > 0 {
			break
		}
		parts := strings.Split(strings.TrimSpace(method), ":")
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		switch k {
		case "header":
			apiKey, err = apiKeyFromHeader(c, v)
		case "query":
			apiKey, err = tokenFromQuery(c, v)
		case "cookie":
			apiKey, err = tokenFromCookie(c, v)
		case "param":
			apiKey, err = tokenFromParam(c, v)
		}
	}
	if err != nil {
		return "", err
	}
	return apiKey, nil
}

// addUserIdToContext adds a given context key and its value to our gin context
func addToContext(c *gin.Context, key consts.ContextKey, value any) *http.Request {
	c.Set(string(key), value)
	return c.Request.WithContext(context.WithValue(c.Request.Context(), key, value))
}

// addUserIdToContext adds user id to our gin context
func addUserIdToContext(c *gin.Context, userID uuid.UUID) *http.Request {
	usrCtxKey := consts.ProjectContextKeys.UserIDCtxKey

	c.Set(string(usrCtxKey), userID)
	return c.Request.WithContext(context.WithValue(c.Request.Context(), usrCtxKey, userID))
}
