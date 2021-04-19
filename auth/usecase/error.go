package usecase

import "net/http"

type baseErr struct {
	error
}
func (_ baseErr) Status() (_ int) { return }
func (_ baseErr) Code() (_ int) { return }

// notFoundErr is error type & used for not found status
type _notFoundErr struct {
	baseErr
}
func notFoundErr(err error) _notFoundErr {
	return _notFoundErr{baseErr{err}}
}
func (_ _notFoundErr) NotFound() {}
func (err _notFoundErr) Status() int { return http.StatusNotFound }

// conflictErr is error type & used for conflict status
type _conflictErr struct {
	baseErr
	code int
}
func conflictErr(err error, code int) _conflictErr {
	return _conflictErr{baseErr{err}, code}
}
func (_ _conflictErr) Conflict() {}
func (err _conflictErr) Status() int { return http.StatusConflict }
func (err _conflictErr) Code() int { return err.code }

// internalServerErr is error type & used for internal server status
type _internalServerErr struct {
	baseErr
}
func internalServerErr(err error) _internalServerErr {
	return _internalServerErr{baseErr{err}}
}
func (_ _internalServerErr) InternalServer() {}
func (err _internalServerErr) Status() int { return http.StatusInternalServerError }
