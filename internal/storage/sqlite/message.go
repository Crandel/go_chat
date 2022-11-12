package sqlite

import (
	"time"

	"github.com/Crandel/go_chat/internal/reading"
)

const MESSAGES = "messages"

type Message struct {
	ID       int       `db:"id,key,auto"`
	Created  time.Time `db:"created"`
	RoomName string    `db:",rel=rooms"`
	UserNick string    `db:",rel=users"`
	Payload  string    `db:"payload"`
}

func (*Message) TableName() string {
	return MESSAGES
}

func (m *Message) ConvertToReading() reading.Message {
	return reading.Message{
		ID:      m.ID,
		Nick:    reading.UserId(m.UserNick),
		Payload: m.Payload,
	}
}
