package memory

import (
	"time"

	r "github.com/Crandel/go_chat/internal/reading"
)

type Message struct {
	Created  time.Time
	UserID   UserId
	RoomName string
	Payload  string
	ID       int
}

func (mm *Message) ConvertMessageToReading() r.Message {
	return r.Message{
		ID:      mm.ID,
		Payload: mm.Payload,
		Nick:    mm.UserID.ConvertUserIdToReading(),
	}
}
