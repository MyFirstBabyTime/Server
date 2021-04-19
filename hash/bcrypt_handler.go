package hash

import "golang.org/x/crypto/bcrypt"

// bcryptHandler is hash handler using bcrypt algorithm
type bcryptHandler struct {}

func BcryptHandler() *bcryptHandler {
	return &bcryptHandler{}
}

// GenerateHashWithMinSalt generate & return hashed value from password with minimum salt
func (bh *bcryptHandler) GenerateHashWithMinSalt(pw string) (string, error) {
	return bh.generateHashFromPW(pw, bcrypt.MinCost)
}

// CompareHashAndPW compare hashed value and password & return error
func (bh *bcryptHandler) CompareHashAndPW(hash, pw string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}

func (bh *bcryptHandler) generateHashFromPW(pw string, salt int) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), salt)
	return string(b), err
}