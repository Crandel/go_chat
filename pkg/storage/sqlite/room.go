package sqlite

import rdn "github.com/Crandel/go_chat/pkg/reading"

const ROOMS = "rooms"

type Room struct {
	ID           string        `db:"id,key"`
	UserMessages []UserMessage `db:"user_messages,rel=user_messages"`
}

func (*Room) TableName() string {
	return ROOMS
}

func (r *Room) ConvertToReading() rdn.Room {
	return rdn.Room{
		Name: r.ID,
	}
}
