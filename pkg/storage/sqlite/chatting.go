package sqlite

import cht "github.com/Crandel/go_chat/pkg/chatting"

func (str *Storage) WriteMessage(u cht.Client, r cht.Room, msg string) error {
	return nil
}

func (str *Storage) ExcludeFromRoom(name string, u cht.Client) error {
	return nil
}

func (str *Storage) AddUserToRoom(name string, u cht.Client) error {
	return nil
}

func (str *Storage) RoomHasUser(name string, u cht.Client) bool {
	return false
}
