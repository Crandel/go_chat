package sqlite

import (
	"fmt"
	"time"

	errs "github.com/Crandel/go_chat/internal/errors"
)

func (str *Storage) AddRoom(rn string) (string, error) {
	const op errs.Op = "sqlite.AddRoom"
	var ret_error error
	res_str := ""
	room := Room{}
	str.db.Select(&room).Where("name = ?", rn).Do()
	if room.Name == rn {
		return res_str, errs.New(op, errs.Info, fmt.Sprintf("Room with name %s already exists", rn))
	}
	room.Name = rn
	room.Created = time.Now()
	error := str.db.Insert(&room).Do()
	if error == nil {
		res_str = room.Name
	} else {
		ret_error = errs.NewError(op, errs.Info, "Failed to create room", error)
	}
	return res_str, ret_error
}
