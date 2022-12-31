package sqlite

import (
	rdn "github.com/Crandel/go_chat/internal/reading"
)

const ROOMS = "rooms"

type Room struct {
	Name string `db:"name,key"`
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
