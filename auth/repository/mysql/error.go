package mysql

// rowNotExistErr is error type & used for row not exist error
type rowNotExistErr struct {
	error
}
func (err rowNotExistErr) Error() string { return err.error.Error() }
func (err rowNotExistErr) RowNotExist() bool { return true }

// isRowNotExist method return if err is row not exist
func isRowNotExist(err error) bool {
	type rowNotExist interface {
		RowNotExist() bool
	}

	re, ok := err.(rowNotExistErr)
	return ok && re.RowNotExist()
}
