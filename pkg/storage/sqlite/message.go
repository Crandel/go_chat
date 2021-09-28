package sqlite

import "time"

const MESSAGES = "messages"

type Message struct {
	ID       int       `db:"id,key,auto"`
	RoomName string    `db:"room_name"`
	UserID   string    `db:"user_id,unique"`
	Payload  string    `db:"payload"`
	Created  time.Time `db:"created"`
}

func (*Message) TableName() string {
	return MESSAGES
}
