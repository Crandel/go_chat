package memory

import "time"

type Message struct {
	ID      string
	Payload string
	Created time.Time
}
