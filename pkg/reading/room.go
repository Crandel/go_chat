package reading

type Room struct {
	Name     string               `json:"name"`
	Messages map[UserId][]Message `json:"messages"`
}
