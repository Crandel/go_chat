package sqlite

import "time"

const MESSAGES = "messages"
const USERMESSAGES = "user_messages"

type Message struct {
	ID      string    `db:"id,key"`
	Payload string    `db:"payload"`
	Created time.Time `db:"created"`
}

func (*Message) TableName() string {
	return MESSAGES
}

type UserMessage struct {
	User     User      `db:"user,rel=users"`
	Messages []Message `db:"messages,rel=messages"`
}

func (*UserMessage) TableName() string {
	return USERMESSAGES
}
