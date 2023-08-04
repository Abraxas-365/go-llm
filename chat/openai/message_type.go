package openai

import (
	"github.com/Abraxas-365/go-llm/chat"
)

func NewMessage(role string, content string) *Message {
	return &Message{Role: role, Content: content}
}

func FromBaseMessage(base chat.BaseMessage) Message {
	return Message{
		Role:    base.GetType(),
		Content: base.GetContent(),
	}
}

func FromBaseMessages(messages []chat.BaseMessage) []Message {
	var result []Message
	for _, base := range messages {
		result = append(result, FromBaseMessage(base))
	}
	return result
}
