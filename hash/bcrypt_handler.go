package hash

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

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
	switch err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)); err {
	case bcrypt.ErrMismatchedHashAndPassword:
		err = mismatchErr{errors.Wrap(err, "failed to CompareHashAndPassword")}
	default:
		err = errors.Wrap(err, "CompareHashAndPassword return unexpected error")
	}
	return
}

func (bh *bcryptHandler) generateHashFromPW(pw string, salt int) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), salt)
	return string(b), err
}

// mismatchErr is error type represent hash & password mismatch error
type mismatchErr struct {
	error
}
func (_ mismatchErr) Mismatch() {}
