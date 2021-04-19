package hash

// bcryptHandler is hash handler using bcrypt algorithm
type bcryptHandler struct {}

func BcryptHandler() *bcryptHandler {
	return &bcryptHandler{}
}
