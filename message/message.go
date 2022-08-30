package message

type Message struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Author  string `json:"author"`
}

type MessageRepository interface {
	CreateMessage(messageContent string, author string) Message
	GetAllFilteredMessages() []Message
}

type MessageResource struct {
	messages []Message
}

func (mr *MessageResource) CreateMessage(messageContent string, author string) {

}

func (mr *MessageResource) GetAllFilteredMessages(count int, offset int) {

}

func NewMessageResource() *MessageResource {
	return &MessageResource{messages: []Message{}}
}
