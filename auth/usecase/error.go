package usecase

// internalServerErr is error type & used for internal server status
type internalServerErr struct {
	error
}
func (_ internalServerErr) InternalServer() {}
