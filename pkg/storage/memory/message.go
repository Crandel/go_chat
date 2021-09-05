package memory

import "time"

type Message struct {
	ID      string
	UserId  UserId
	Payload string
	Created time.Time
}
