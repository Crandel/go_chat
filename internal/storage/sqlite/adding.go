package sqlite

import (
	"fmt"

	lg "github.com/Crandel/go_chat/internal/logging"
)

func (str *Storage) AddRoom(rn string) (string, error) {
	const op lg.Stk = "sqlite.AddRoom"
	if rn == "" {
		return "", lg.New(
			op, "Room name could not be empty",
		)
	}

	res_str := ""
	room := Room{}
	error := str.db.Select(&room).Where("name = ?", rn).Do()
	if error == nil && room.Name == rn {
		return "", lg.New(op, fmt.Sprintf("Room with name %s already exists", rn))
	}
	room.Name = rn
	error = str.db.Insert(&room).Do()
	if error == nil {
		res_str = room.Name
	} else {
		return "", lg.NewError(op, "Failed to create room", error)
	}
	return res_str, nil
}
