package memory

import (
	"fmt"
	"time"

	errs "github.com/Crandel/go_chat/pkg/errors"
)

func (str *Storage) AddRoom(rn string) (string, []error) {
	const op errs.Op = "memory.AddRoom"
	var error_list []error
	str.Lock()
	if str.Rooms == nil {
		str.Rooms = make(map[string]Room)
	} else {
		_, exists := str.Rooms[rn]
		if exists {
			error_list = append(
				error_list,
				errs.New(
					op, errs.Info, fmt.Sprintf("Room with name %s already exists", rn)))
			return "", error_list
		}
	}
	str.Unlock()
	mr := Room{Name: rn, Created: time.Now()}
	str.Lock()
	str.Rooms[mr.Name] = mr
	str.Unlock()
	return mr.Name, error_list
}
