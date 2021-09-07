package sqlite

type Room struct {
	ID       string      `db:"id,key"`
	Messages UserMessage `db:"messages"`
}
