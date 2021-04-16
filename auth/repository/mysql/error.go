package mysql

// rowNotExistErr is error type & used for row not exist error
type rowNotExistErr struct {
	error
}
func (err rowNotExistErr) Error() string { return err.error.Error() }
func (err rowNotExistErr) IsRowNotExist() bool { return true }

// rowNotExist interface & isRowNotExist method is used for check & get error context
type rowNotExist interface {
	IsRowNotExist() bool
}
func isRowNotExist(err error) (bool, rowNotExist) {
	re, ok := err.(rowNotExistErr)
	return ok && re.IsRowNotExist(), re
}

// entryDuplicateErr is error type & used for row not exist error
type entryDuplicateErr struct {
	error
	duplicateKey string
}
func (err entryDuplicateErr) Error() string { return err.error.Error() }
func (err entryDuplicateErr) IsEntryDuplicate() bool { return true }
func (err entryDuplicateErr) DuplicateKey() string { return err.duplicateKey }

// entryDuplicate interface & isEntryDuplicate method is used for check & get error context
type entryDuplicate interface {
	IsEntryDuplicate() bool
	DuplicateKey() string
}
func isEntryDuplicate(err error) (bool, entryDuplicate) {
	ee, ok := err.(entryDuplicate)
	return ok && ee.IsEntryDuplicate(), ee
}
