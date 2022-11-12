package sqlite

import (
	"github.com/samonzeweb/godb"
)

type Storage struct {
	db *godb.DB
}

func NewStorage(db *godb.DB) Storage {
	return Storage{db}
}
