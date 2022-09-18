package reading

type Message struct {
	Nick    UserId `json:"nick"`
	Payload string `json:"payload"`
	ID      int    `json:"id"`
}
