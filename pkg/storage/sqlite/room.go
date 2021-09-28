package sqlite

import (
	"time"

	rdn "github.com/Crandel/go_chat/pkg/reading"
)

const ROOMS = "rooms"

type Room struct {
	Name         string        `db:"name,key"`
	UserMessages []UserMessage `db:"user_messages,rel=user_messages"`
	Created      time.Time
}

func (*Room) TableName() string {
	return ROOMS
}

func (r *Room) ConvertToReading() rdn.Room {
	return rdn.Room{
		Name:     r.Name,
		Messages: map[rdn.UserId][]rdn.Message{},
	}
}
