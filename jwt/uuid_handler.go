package jwt

// uuidHandler is jwt handler about uuid token
type uuidHandler struct {
	jwtKey string
}

func UUIDHandler(key string) *uuidHandler {
	return &uuidHandler{
		jwtKey: key,
	}
}
