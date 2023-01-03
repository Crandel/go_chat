package memory

import (
	lg "github.com/Crandel/go_chat/internal/logging"
	rdn "github.com/Crandel/go_chat/internal/reading"
)

func (str *Storage) ReadRooms() ([]rdn.Room, error) {
	var rooms = []rdn.Room{}
	str.RLock()
	for _, room := range str.Rooms {
		rooms = append(rooms, str.collectRoomMessages(room.Name))
	}
	str.RUnlock()
	return rooms, nil
}

func (str *Storage) ReadRoom(rid string) (rdn.Room, error) {
	const op lg.Op = "memory.ReadRoom"
	str.RLock()
	room, exists := str.Rooms[rid]
	str.RUnlock()
	if !exists {
		return rdn.Room{}, lg.New(op, "No rooms with id "+rid)
	}
	return str.collectRoomMessages(room.Name), nil
}

func (str *Storage) ReadUsers() ([]rdn.User, error) {
	const op lg.Op = "memory.ReadUsers"
	var users = []rdn.User{}
	str.RLock()
	for _, u := range str.Users {
		users = append(users, u.ConvertUserToReading())
	}
	str.RUnlock()
	var not_found error
	if len(users) == 0 {
		not_found = lg.New(op, "No users are here")
	}
	return users, not_found
}

func (str *Storage) ReadUser(uid rdn.UserId) (rdn.User, error) {
	const op lg.Op = "memory.ReadUser"
	umid := ConvertUserIdFromReading(uid)
	str.RLock()
	s_user, exists := str.Users[umid]
	str.RUnlock()
	if !exists {
		return rdn.User{}, lg.New(op, "No user with id "+string(uid))
	}
	return s_user.ConvertUserToReading(), nil
}

func (str *Storage) collectRoomMessages(name string) rdn.Room {
	rRoom := rdn.Room{Name: name}
	rMessages := make(map[rdn.UserId][]rdn.Message)
	str.RLock()
	for _, m := range str.Messages {
		if m.RoomName == name {
			ruid := m.UserID.ConvertUserIdToReading()
			rMessage := rdn.Message{
				ID:      m.ID,
				Nick:    ruid,
				Payload: m.Payload,
			}
			rMessages[ruid] = append(rMessages[ruid], rMessage)
		}
	}
	str.RUnlock()
	if len(rMessages) > 0 {
		rRoom.Messages = rMessages
	}
	return rRoom
}

// func (sr *Room) ConvertRoomToReading(messages []Message) rdn.Room {
// 	var m_messages = map[rdn.UserId][]rdn.Message{}
// 	for _, s_messages := range messages {
// 		var l_messages = []rdn.Message{}
// 		for _, s_message := range s_messages {
// 			l_messages = append(l_messages, s_message.ConvertMessageToReading())
// 		}
// 		m_messages[s_ui.ConvertUserIdToReading()] = l_messages
// 	}
// 	return r.Room{Name: sr.Name, Messages: m_messages}
// }
