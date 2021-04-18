package usecase

// rowNotExist interface & isRowNotExist method is used for check & get error context
type rowNotExistErr interface {
	error
	RowNotExist()
}

// entryDuplicate interface & isEntryDuplicate method is used for check & get error context
type entryDuplicateErr interface {
	error
	EntryDuplicate()
	DuplicateKey() string
}

// noReferencedRow interface & isNoReferencedRow method is used for check & get error context
type noReferencedRowErr interface {
	error
	NoReferencedRow()
	ForeignKey() string
}

// invalidModel interface & isInvalidModel method is used for check & get error context
type invalidModelErr interface {
	error
	InvalidModel()
}
