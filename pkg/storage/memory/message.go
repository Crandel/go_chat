package memory

import "time"

type Message struct {
	ID      string
	Text    string
	Created time.Time
}
