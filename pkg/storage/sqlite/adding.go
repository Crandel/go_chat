package sqlite

import (
	"time"

	add "github.com/Crandel/go_chat/pkg/adding"
	errs "github.com/Crandel/go_chat/pkg/errors"
)

func (str *Storage) AddRoom(ar add.Room) (string, []error) {
	const op errs.Op = "sqlite.AddRoom"
	room := Room{}
	str.db.Select(&room).Where("name = ?", ar.Name).Do()
	if room.Name == ar.Name {
		return "", []error{errs.New(op, errs.Info, "Room already exists")}
	}
	list_errors := []error{}
	room.Name = ar.Name
	room.Created = time.Now()
	error := str.db.Insert(&room).Do()
	res_str := ""
	if error == nil {
		res_str = room.Name
	} else {
		list_errors = append(
			list_errors,
			errs.NewError(op, errs.Info, "Failed to create room", error))
	}
	return res_str, list_errors
}
