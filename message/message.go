package message

type Message struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Author  string `json:"author"`
}

type MessageRepository interface {
	CreateMessage(messageContent string, author string) Message
	GetPaginatedMessages(count int, offset int) []Message
}

type MessageResource struct {
	messages []Message
}

func (mr *MessageResource) CreateMessage(messageContent string, author string) Message {
	message := Message{ID: string(len(mr.messages)), Author: author, Message: messageContent}
	mr.messages = append(mr.messages, message)
	return message
}

func (mr *MessageResource) GetPaginatedMessages(count int, offset int) []Message {
	if count > 100 {
		count = 100
	}
	if count < 1 {
		count = 10
	}
	if offset < 0 {
		offset = 0
	}
	if offset > len(mr.messages) {
		offset = len(mr.messages)
	}
	if offset+count > len(mr.messages) {
		return mr.messages[offset:]
	}
	return mr.messages[offset : offset+count]
}

func NewMessageResource() *MessageResource {
	return &MessageResource{messages: []Message{}}
}
