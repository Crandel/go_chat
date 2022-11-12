package memory

import (
	"time"
)

type Room struct {
	Created time.Time
	Name    string
	Members []UserId
}
