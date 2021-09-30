package sqlite

import (
	"time"

	rdn "github.com/Crandel/go_chat/pkg/reading"
)

const ROOMS = "rooms"

type Room struct {
	Name    string    `db:"name,key"`
	Created time.Time `db:"created"`
}

func (*Room) TableName() string {
	return ROOMS
}

func (r *Room) ConvertToReading(messages map[rdn.UserId][]rdn.Message) rdn.Room {
	return rdn.Room{
		Name:     r.Name,
		Messages: messages,
	}
}
