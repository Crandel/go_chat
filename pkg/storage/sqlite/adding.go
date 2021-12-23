package sqlite

import (
	"time"

	errs "github.com/Crandel/go_chat/pkg/errors"
)

func (str *Storage) AddRoom(rn string) (string, []error) {
	const op errs.Op = "sqlite.AddRoom"
	room := Room{}
	str.db.Select(&room).Where("name = ?", rn).Do()
	if room.Name == rn {
		return "", []error{errs.New(op, errs.Info, "Room already exists")}
	}
	list_errors := []error{}
	room.Name = rn
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
