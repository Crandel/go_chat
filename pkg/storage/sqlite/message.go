package sqlite

import "time"

type Message struct {
	ID      string    `db:"id,key"`
	Payload string    `db:"payload"`
	Created time.Time `db:"created"`
}

type UserMessage struct {
	Message Message `db:"message"`
	User    User    `db:"user"`
}
