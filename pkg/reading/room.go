package reading

type Room struct {
	Name     string
	Messages map[UserId][]Message
}
