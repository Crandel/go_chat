package memory

import (
	"fmt"
	"time"

	lg "github.com/Crandel/go_chat/internal/logging"
)

func (str *Storage) AddRoom(rn string) (string, error) {
	const op lg.Op = "memory.AddRoom"
	str.Lock()
	if str.Rooms == nil {
		str.Rooms = make(map[string]Room)
	} else {
		_, exists := str.Rooms[rn]
		if exists {
			return "", lg.New(
				op, lg.Info, fmt.Sprintf("Room with name %s already exists", rn))
		}
	}
	str.Unlock()
	mr := Room{Name: rn, Created: time.Now()}
	str.Lock()
	str.Rooms[mr.Name] = mr
	str.Unlock()
	return mr.Name, nil
}
