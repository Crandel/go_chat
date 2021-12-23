package memory

import (
	"time"
)

type Room struct {
	Name    string
	Members []UserId
	Created time.Time
}
