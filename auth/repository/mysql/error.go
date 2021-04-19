package mysql

// rowNotExistErr is error type & used for row not exist error
type rowNotExistErr struct {
	error
}
func (_ rowNotExistErr) RowNotExist() {}

// entryDuplicateErr is error type & used for row not exist error
type entryDuplicateErr struct {
	error
	duplicateKey string
}
func (_ entryDuplicateErr) EntryDuplicate() {}
func (err entryDuplicateErr) DuplicateKey() string { return err.duplicateKey }

// noReferencedRowErr is error type & used for no referenced row error
type noReferencedRowErr struct {
	error
	foreignKey string
}
func (_ noReferencedRowErr) NoReferencedRow() {}
func (err noReferencedRowErr) ForeignKey() string { return err.foreignKey }

// invalidModelErr is error type & used for row not exist error
type invalidModelErr struct {
	error
}
func (_ invalidModelErr) InvalidModel() {}
