package sqlite

import (
	"time"
)

const ROOMS = "rooms"

type Room struct {
	Name    string    `db:"name,key"`
	Created time.Time `db:"created"`
}

func (*Room) TableName() string {
	return ROOMS
}
