package sqlite

import (
	"time"

	"github.com/Crandel/go_chat/internal/reading"
)

const MESSAGES = "messages"

type Message struct {
	ID         int       `db:"id,key,auto"`
	Created    time.Time `db:"created"`
	UserRoomID int       `db:"user_room_id"`
	Payload    string    `db:"payload"`
}

func (*Message) TableName() string {
	return MESSAGES
}

func (m *Message) ConvertToReading(nick string) reading.Message {
	return reading.Message{
		ID:      m.ID,
		Nick:    reading.UserId(nick),
		Payload: m.Payload,
	}
}
