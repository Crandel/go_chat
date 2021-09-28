package reading

type Message struct {
	ID      int    `json:"id"`
	UserId  UserId `json:"user_id"`
	Payload string `json:"payload"`
}
