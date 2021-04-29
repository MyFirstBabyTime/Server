package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
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
func (uh *uuidHandler) ParseUUIDFromToken(s string) (uuid, _type string, err error) {
	token, err := jwt.ParseWithClaims(s, &uuidClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(uh.jwtKey), nil
	})
	if err != nil {
		err = errors.Wrap(err, "failed to ParseWithClaims")
		return
	}

	claims, ok := token.Claims.(*uuidClaims)
	if !ok || !token.Valid {
		err = errors.Wrap(err, "failed to assert to *uuidClaims")
		return
	}

	uuid = claims.UUID
	_type = claims.Type
	return
}
