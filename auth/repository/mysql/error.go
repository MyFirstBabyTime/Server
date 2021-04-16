package mysql

// rowNotExistErr is error type & used for row not exist error
type rowNotExistErr struct {
	error
}
func (err rowNotExistErr) Error() string { return err.error.Error() }
func (err rowNotExistErr) IsRowNotExist() bool { return true }

// entryDuplicateErr is error type & used for row not exist error
type entryDuplicateErr struct {
	error
	duplicateEntry string
}
func (err entryDuplicateErr) Error() string { return err.error.Error() }
func (err entryDuplicateErr) IsEntryDuplicate() bool { return true }
func (err entryDuplicateErr) DuplicateEntry() string { return err.duplicateEntry }

// isRowNotExist method return if err is about row not exist
func isRowNotExist(err error) bool {
	type rowNotExist interface {
		IsRowNotExist() bool
	}

	re, ok := err.(rowNotExistErr)
	return ok && re.IsRowNotExist()
}

// isEntryDuplicate method return if err is about entry duplicate
func isEntryDuplicate(err error) (bool, string) {
	type entryDuplicate interface {
		IsEntryDuplicate() bool
		DuplicateEntry() string
	}

	re, ok := err.(entryDuplicate)
	return ok && re.IsEntryDuplicate(), re.DuplicateEntry()
}
