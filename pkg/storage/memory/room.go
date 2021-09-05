package memory

import (
	"errors"
	"fmt"
	"time"

	"github.com/Crandel/go_chat/pkg/adding"
)

type Room struct {
	Name     string
	Messages map[UserId][]Message

	Created time.Time
}

func (s *Storage) AddRoom(ar adding.Room) (string, error) {
	_, exists := s.Rooms[ar.Name]
	if exists {
		return "", errors.New(fmt.Sprintf("Room with name %s already exists", ar.Name))
	}
	messages := make(map[UserId][]Message)
	for _, user := range ar.Users {
		messages[UserId(user.ID)] = []Message{}
	}
	mr := Room{Name: ar.Name, Messages: messages, Created: time.Now()}
	s.Rooms[mr.Name] = mr
	return mr.Name, nil
}

func (s *Storage) ReadRooms() ([]Room, error) {
	var rooms = []Room{}
	for _, room := range s.Rooms {
		var map_messages = map[UserId][]Message{}
		for s_ui, s_messages := range room.Messages {
			var m_messages = []Message{}
			for _, s_message := range s_messages {
				m_message := Message{
					ID:      s_message.ID,
					Payload: s_message.Payload,
					UserId:  UserId(s_ui),
				}
				m_messages = append(m_messages, m_message)
			}
			map_messages[s_ui] = m_messages
		}
		rooms = append(rooms, Room{
			Name:     room.Name,
			Messages: map_messages,
		})
	}
	return rooms, nil
}
