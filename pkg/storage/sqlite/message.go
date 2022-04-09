package sqlite

import (
	"time"

	"github.com/Crandel/go_chat/pkg/reading"
)

const MESSAGES = "messages"

type Message struct {
	Created  time.Time `db:"created"`
	RoomName string    `db:"room_name"`
	UserID   string    `db:"user_id,unique"`
	Payload  string    `db:"payload"`
	ID       int       `db:"id,key,auto"`
}

func (*Message) TableName() string {
	return MESSAGES
}

func (m *Message) ConvertToReading() reading.Message {
	return reading.Message{
		ID:      m.ID,
		UserId:  reading.UserId(m.UserID),
		Payload: m.Payload,
	}
}
