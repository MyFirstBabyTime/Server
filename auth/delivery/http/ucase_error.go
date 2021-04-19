package http

// usecaseErr is interface type to assert error defined in usecase
type usecaseErr interface {
	error
	Status() int
	Code() int
}
