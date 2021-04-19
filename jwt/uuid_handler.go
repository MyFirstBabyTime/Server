package jwt

import (
	"github.com/dgrijalva/jwt-go"
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
