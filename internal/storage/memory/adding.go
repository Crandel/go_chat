package memory

import (
	"fmt"
	"time"

	errs "github.com/Crandel/go_chat/internal/errors"
)

func (str *Storage) AddRoom(rn string) (string, error) {
	const op errs.Op = "memory.AddRoom"
	str.Lock()
	if str.Rooms == nil {
		str.Rooms = make(map[string]Room)
	} else {
		_, exists := str.Rooms[rn]
		if exists {
			return "", errs.New(
				op, errs.Info, fmt.Sprintf("Room with name %s already exists", rn))
		}
	}
	str.Unlock()
	mr := Room{Name: rn, Created: time.Now()}
	str.Lock()
	str.Rooms[mr.Name] = mr
	str.Unlock()
	return mr.Name, nil
}
