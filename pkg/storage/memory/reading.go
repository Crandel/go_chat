package memory

import (
	errs "github.com/Crandel/go_chat/pkg/errors"
	rdn "github.com/Crandel/go_chat/pkg/reading"
)

func (str *Storage) ReadRooms() ([]rdn.Room, error) {
	const op errs.Op = "memory.ReadRooms"
	var rooms = []rdn.Room{}
	str.RLock()
	for _, room := range str.Rooms {
		rooms = append(rooms, str.collectRoomMessages(room.Name))
	}
	str.RUnlock()
	return rooms, nil
}

func (str *Storage) ReadRoom(rid string) (rdn.Room, error) {
	const op errs.Op = "memory.ReadRoom"
	str.RLock()
	room, exists := str.Rooms[rid]
	str.RUnlock()
	if !exists {
		return rdn.Room{}, errs.New(op, errs.Info, "No rooms with id "+rid)
	}
	return str.collectRoomMessages(room.Name), nil
}

func (str *Storage) ReadUsers() ([]rdn.User, error) {
	const op errs.Op = "memory.ReadUsers"
	var users = []rdn.User{}
	str.RLock()
	for _, u := range str.Users {
		users = append(users, u.ConvertUserToReading())
	}
	str.RUnlock()
	var not_found error
	if len(users) == 0 {
		not_found = errs.New(op, errs.Info, "No users are here")
	}
	return users, not_found
}

func (str *Storage) ReadUser(uid rdn.UserId) (rdn.User, error) {
	const op errs.Op = "memory.ReadUser"
	umid := ConvertUserIdFromReading(uid)
	str.RLock()
	s_user, exists := str.Users[umid]
	str.RUnlock()
	if !exists {
		return rdn.User{}, errs.New(op, errs.Info, "No user with id "+string(uid))
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
			innerMessages, _ := rMessages[ruid]
			innerMessages = append(innerMessages, rMessage)
			rMessages[ruid] = innerMessages
		}
	}
	str.RUnlock()
	if len(rMessages) > 0 {
		rRoom.Messages = rMessages
	}
	return rRoom
}

// func (sr *Room) ConvertRoomToReading() r.Room {
// 	var m_messages = map[r.UserId][]r.Message{}
// 	for s_ui, s_messages := range sr.Messages {
// 		var l_messages = []r.Message{}
// 		for _, s_message := range s_messages {
// 			l_messages = append(l_messages, s_message.ConvertMessageToReading())
// 		}
// 		m_messages[s_ui.ConvertUserIdToReading()] = l_messages
// 	}
// 	return r.Room{Name: sr.Name, Messages: m_messages}
// }
