package sqlite

const ROOMS = "rooms"

type Room struct {
	ID           string        `db:"id,key"`
	UserMessages []UserMessage `db:"user_messages,rel=user_messages"`
}

func (*Room) TableName() string {
	return ROOMS
}
