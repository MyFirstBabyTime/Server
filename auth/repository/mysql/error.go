package mysql

// rowNotExistErr is error type & used for row not exist error
type rowNotExistErr struct {
	error
}
func (err rowNotExistErr) Error() string { return err.error.Error() }
func (err rowNotExistErr) RowNotExist() bool { return true }