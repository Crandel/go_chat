package sqlite

import (
	"fmt"
	"time"

	lg "github.com/Crandel/go_chat/internal/logging"
)

func (str *Storage) AddRoom(rn string) (string, error) {
	const op lg.Op = "sqlite.AddRoom"
	res_str := ""
	room := Room{}
	error := str.db.Select(&room).Where("name = ?", rn).Do()
	if error != nil || room.Name == rn {
		return "", lg.New(op, lg.Info, fmt.Sprintf("Room with name %s already exists", rn))
	}
	room.Name = rn
	room.Created = time.Now()
	error = str.db.Insert(&room).Do()
	if error == nil {
		res_str = room.Name
	} else {
		return "", lg.NewError(op, lg.Info, "Failed to create room", error)
	}
	return res_str, nil
}
