package reading

type Message struct {
	ID      string `json:"id"`
	UserId  UserId `json:"user_id"`
	Payload string `json:"payload"`
}
