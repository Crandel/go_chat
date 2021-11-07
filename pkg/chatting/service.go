package chatting

var commands chan command

type Repository interface {
	AddMessage(Message)
}
