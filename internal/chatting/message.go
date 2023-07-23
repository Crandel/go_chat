package chatting

type ChatMessage struct {
	CommandId CommandID `json:"command"`
	User      *string   `json:"username"`
	Args      []string  `json:"args"`
}
