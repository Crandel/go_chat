package memory

import (
	"time"

	lg "github.com/Crandel/go_chat/internal/logging"
)

const EmptyRoomNameError = "Room name could not be empty"

func (str *Storage) AddRoom(rn string) (string, error) {
	const op lg.Stk = "memory.AddRoom"
	if rn == "" {
		return "", lg.New(
			op, EmptyRoomNameError,
		)
	}
	str.Lock()
	if str.Rooms == nil {
		str.Rooms = make(map[string]Room)
	} else {
		_, exists := str.Rooms[rn]
		if exists {
			return rn, nil
		}
	}
	str.Unlock()
	mr := Room{Name: rn, Created: time.Now()}
	str.Lock()
	str.Rooms[mr.Name] = mr
	str.Unlock()
	return mr.Name, nil
}
