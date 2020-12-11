package domainErr

import "errors"

var (
	DuplicateKeyError = errors.New("Duplicate key value")
)
