package reading

type Room struct {
	Messages map[UserId][]Message `json:"messages"`
	Name     string               `json:"name"`
}
