package domainErr

import "errors"

var (
	DuplicateKeyError = errors.New("Duplicate key value")
	NotExists = errors.New("Item not exists")
	ThreadNotExists = errors.New("Thread not exists")
	PostNotExists = errors.New("Post not exists")
	AlreadyExists = errors.New("Already exists exists")
	EmptyReq = errors.New("Empty Request")
)
