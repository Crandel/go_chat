package memory

import cht "github.com/Crandel/go_chat/pkg/chatting"

func (str *Storage) WriteMessage(u cht.User, r cht.Room, msg string) error {
	return nil
}

func (str *Storage) ExcludeFromRoom(roomName string, u cht.User) error {
	return nil
}

func (str *Storage) AddUserToRoom(roomName string, u cht.User) error {

	return nil
}

func (str *Storage) RoomHasUser(roomName string, u cht.User) bool {

	return false
}
