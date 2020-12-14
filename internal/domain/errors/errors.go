package domainErr

import "errors"

var (
	DuplicateKeyError = errors.New("Duplicate key value")
	NotExists = errors.New("Item not exists")
)
