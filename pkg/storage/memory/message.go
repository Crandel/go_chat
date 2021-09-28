package memory

import (
	"time"

	r "github.com/Crandel/go_chat/pkg/reading"
)

type Message struct {
	ID       int
	UserId   UserId
	RoomName string
	Payload  string
	Created  time.Time
}

func (mm *Message) ConvertMessageToReading() r.Message {
	return r.Message{
		ID:      mm.ID,
		Payload: mm.Payload,
		UserId:  mm.UserId.ConvertUserIdToReading(),
	}
}
