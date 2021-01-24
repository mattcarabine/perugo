package db

import "errors"

// Entities
// Users
// Rooms
// Games

type DB interface {
	LookupId(string, *interface{}) error
	Store(string, *interface{}) error
}

var (
	ErrDBEntityDoesNotExist = errors.New("entity does not exist")
	ErrDBLookupFailed       = errors.New("lookup failed due to internal issue")
	ErrDBStoreFailed        = errors.New("store operation failed due to internal issue")
)
