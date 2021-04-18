package usecase

import "net/http"

// notFoundErr is error type & used for not found status
type notFoundErr struct {
	error
}
func (_ notFoundErr) NotFound() {}
func (err notFoundErr) Status() int { return http.StatusNotFound }

// conflictErr is error type & used for conflict status
type conflictErr struct {
	error
	code int
}
func (_ conflictErr) Conflict() {}
func (err conflictErr) Status() int { return http.StatusConflict }
func (err conflictErr) Code() int { return err.code }

// internalServerErr is error type & used for internal server status
type internalServerErr struct {
	error
}
func (_ internalServerErr) InternalServer() {}
func (err internalServerErr) Status() int { return http.StatusInternalServerError }
