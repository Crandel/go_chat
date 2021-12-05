package memory

import (
	"fmt"
	"time"

	add "github.com/Crandel/go_chat/pkg/adding"
	errs "github.com/Crandel/go_chat/pkg/errors"
)

func (str *Storage) AddRoom(ar add.Room) (string, []error) {
	const op errs.Op = "memory.AddRoom"
	var error_list []error
	str.Lock()
	if str.Rooms == nil {
		str.Rooms = make(map[string]Room)
	} else {
		_, exists := str.Rooms[ar.Name]
		if exists {
			error_list = append(
				error_list,
				errs.New(
					op, errs.Info, fmt.Sprintf("Room with name %s already exists", ar.Name)))
			return "", error_list
		}
	}
	str.Unlock()
	mr := Room{Name: ar.Name, Created: time.Now()}
	str.Lock()
	str.Rooms[mr.Name] = mr
	str.Unlock()
	return mr.Name, error_list
}
