package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type DataAccessObjects struct {
	Movies MovieDAO
}

func NewDataAccessObjects(db *sql.DB) DataAccessObjects {
	return DataAccessObjects{
		Movies: MovieDAO{DB: db},
	}
}
