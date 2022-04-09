package reading

type Message struct {
	UserId  UserId `json:"user_id"`
	Payload string `json:"payload"`
	ID      int    `json:"id"`
}
