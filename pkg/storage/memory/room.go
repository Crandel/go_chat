package memory

import (
	"time"

	r "github.com/Crandel/go_chat/pkg/reading"
)

type Room struct {
	Name     string
	Messages map[UserId][]Message

	Created time.Time
}

func (sr *Room) ConvertRoomToReading() r.Room {
	var m_messages = map[r.UserId][]r.Message{}
	for s_ui, s_messages := range sr.Messages {
		var l_messages = []r.Message{}
		for _, s_message := range s_messages {
			l_messages = append(l_messages, s_message.ConvertMessageToReading())
		}
		m_messages[s_ui.ConvertUserIdToReading()] = l_messages
	}
	return r.Room{Name: sr.Name, Messages: m_messages}
}
