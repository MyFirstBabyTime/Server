package jwt

import (
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

// ParseUUIDFromToken is middleware that parse uuid & type from token received from request header
func (uh *uuidHandler) ParseUUIDFromToken(c *gin.Context) {
	var tokenStr string
	if tokens := c.Request.Header["Authorization"]; len(tokens) >= 1 {
		tokenStr = tokens[0]
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, defaultResp(http.StatusUnauthorized, 0, "Authorization not set"))
		return
	}

	token, err := jwt.ParseWithClaims(tokenStr, &uuidClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(uh.jwtKey), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, defaultResp(http.StatusUnauthorized, 0, err.Error()))
		return
	}

	claims, ok := token.Claims.(*uuidClaims)
	if !ok || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, defaultResp(http.StatusUnauthorized, 0, "failed to assert token Claim"))
		return
	}

	c.Set("uuid", claims.UUID)
	c.Set("_type", claims.Type)
	c.Next() // middleware로 쓰인다는 것을 명시하기 위해 c.Next() 호출 (호출 안해도 다음으로 등록된 handler 실행되긴 함)
}

// defaultResp return response have status, code, message inform
func defaultResp(status, code int, msg string) (resp gin.H) {
	resp = gin.H{}
	resp["status"] = status
	resp["code"] = code
	resp["message"] = msg
	return
}
