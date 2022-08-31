package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessageResource_CreateMessage(t *testing.T) {
	mr := NewMessageResource()
	messages := mr.GetPaginatedMessages(1, 0)
	assert.Equal(t, len(messages), 0)
	mr.CreateMessage("hello world", "swong")
	messages2 := mr.GetPaginatedMessages(1, 0)
	assert.Equal(t, len(messages2), 1)
}

func TestMessageResource_GetPaginatedMessages(t *testing.T) {
	mr := NewMessageResource()
	mr.CreateMessage("hello world", "swong")
	mr.CreateMessage("hello world", "swong")
	mr.CreateMessage("hello world", "swong")
	messages := mr.GetPaginatedMessages(3, 1)
	assert.Equal(t, len(messages), 2)
	messages2 := mr.GetPaginatedMessages(100, 0)
	assert.Equal(t, len(messages2), 3)
	messages3 := mr.GetPaginatedMessages(3, 100)
	assert.Equal(t, len(messages3), 1)
	messages4 := mr.GetPaginatedMessages(1, 0)
	assert.Equal(t, len(messages4), 1)
}
