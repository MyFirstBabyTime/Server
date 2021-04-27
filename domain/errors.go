package domain

//
type RepoErr error

// ErrRowNotExistErr is error type & used for row not exist error
type ErrRowNotExist struct {
	RepoErr
}

// ErrEntryDuplicate is error type & used for row not exist error
type ErrEntryDuplicate struct {
	RepoErr
	DuplicateKey string
}

// ErrNoReferencedRow is error type & used for no referenced row error
type ErrNoReferencedRow struct {
	RepoErr
	ForeignKey string
}

// ErrInvalidModel is error type & used for row not exist error
type ErrInvalidModel struct {
	RepoErr
}

//
type UsecaseErr error

//
type UsecaseError struct {
	UsecaseErr
	Status, Code int
}
