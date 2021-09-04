package memory

import "time"

type Room struct {
	ID    string
	Users []User

	Created time.Time
}
