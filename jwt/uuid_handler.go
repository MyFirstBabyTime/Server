package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// uuidHandler is jwt handler about uuid token
type uuidHandler struct {
	jwtKey string
}

func UUIDHandler(key string) *uuidHandler {
	return &uuidHandler{
		jwtKey: key,
	}
}

// uuidClaims is used for generate JWT including uuid inform
type uuidClaims struct {
	UUID string `json:"uuid"`
	Type string `json:"type"`
	jwt.StandardClaims
}

// GenerateUUIDJWT generate & return JWT UUID token with type & time
func (uh *uuidHandler) GenerateUUIDJWT(uuid, _type string, t time.Duration) (token string, err error) {
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS512, uuidClaims{
		UUID: uuid,
		Type: _type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(t).Unix(),
		},
	}).SignedString([]byte(uh.jwtKey))
	return
}

// ParseUUIDFromToken parse uuid & type from token received from parameter
func (uh *uuidHandler) ParseUUIDFromToken (c *gin.Context) {
	accessToken := c.Request.Header["Authorization"][0]

	token, err := jwt.ParseWithClaims(accessToken, &uuidClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(uh.jwtKey), nil
	})
	fmt.Println(token, err)

	if err != nil {
		c.JSON(http.StatusUnauthorized, defaultResp(http.StatusUnauthorized, 0, err.Error()))
		c.Abort()
		return
	}

	claims, ok := token.Claims.(*uuidClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, defaultResp(http.StatusUnauthorized, 0, err.Error()))
		c.Abort()
		return
	}

	uuid = claims.UUID
	_type = claims.Type
	return
}
