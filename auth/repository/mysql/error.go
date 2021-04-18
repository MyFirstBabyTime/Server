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

// noReferencedRowErr is error type & used for no referenced row error
type noReferencedRowErr struct {
	error
	foreignKey string
}
func (err noReferencedRowErr) Error() string { return err.error.Error() }
func (err noReferencedRowErr) IsNoReferenced() bool { return true }
func (err noReferencedRowErr) ForeignKey() string { return err.foreignKey }

// noReferencedRow interface & isNoReferencedRow method is used for check & get error context
type noReferencedRow interface {
	IsNoReferencedRow() bool
	ForeignKey() string
}
func isNoReferencedRow(err error) (bool, noReferencedRow) {
	nr, ok := err.(noReferencedRow)
	return ok && nr.IsNoReferencedRow(), nr
}

// invalidModelErr is error type & used for row not exist error
type invalidModelErr struct {
	error
}
func (err invalidModelErr) Error() string { return err.error.Error() }
func (err invalidModelErr) IsInvalidModel() bool { return true }

// invalidModel interface & isInvalidModel method is used for check & get error context
type invalidModel interface {
	IsInvalidModel() bool
}
func isInvalidModel(err error) (bool, invalidModel) {
	im, ok := err.(invalidModel)
	return ok && im.IsInvalidModel(), im
}
