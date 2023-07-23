package chatting

type ChatMessage struct {
	CommandId CommandID `json:"command"`
	Args      []string  `json:"args"`
}
