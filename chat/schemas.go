package chat

import (
	"errors"
)

type BaseMessage interface {
	GetType() string
	GetContent() string
}

type HumanMessage struct {
	content string
}

func NewHumanMessage(content string) BaseMessage {
	return &HumanMessage{content: content}
}

func (h *HumanMessage) GetType() string {
	return "user"
}

func (h *HumanMessage) GetContent() string {
	return h.content
}

type SystemMessage struct {
	content string
}

func NewSystemMessage(content string) BaseMessage {
	return &SystemMessage{content: content}
}

func (s *SystemMessage) GetType() string {
	return "system"
}

func (s *SystemMessage) GetContent() string {
	return s.content
}

type AIMessage struct {
	content string
}

func NewAIMessage(content string) BaseMessage {
	return &AIMessage{content: content}
}

func (a *AIMessage) GetType() string {
	return "assistant"
}

func (a *AIMessage) GetContent() string {
	return a.content
}

func MessageFromMap(message map[string]string) (BaseMessage, error) {
	messageType, ok := message["type"]
	if !ok {
		return nil, errors.New("No type key on map")
	}

	content, ok := message["content"]
	if !ok {
		content = ""
	}

	switch messageType {
	case "user":
		return NewHumanMessage(content), nil
	case "system":
		return NewSystemMessage(content), nil
	case "assistant":
		return NewAIMessage(content), nil
	default:
		return nil, errors.New("Got unexpected message type: " + messageType)
	}
}

func MessagesFromMap(messages []map[string]string) ([]BaseMessage, error) {
	var result []BaseMessage
	for _, message := range messages {
		msg, err := MessageFromMap(message)
		if err != nil {
			return nil, err
		}
		result = append(result, msg)
	}
	return result, nil
}

func MessageToMap(message BaseMessage) map[string]string {
	return map[string]string{
		"type":    message.GetType(),
		"content": message.GetContent(),
	}
}

func MessagesToMap(messages []BaseMessage) []map[string]string {
	var result []map[string]string
	for _, message := range messages {
		result = append(result, MessageToMap(message))
	}
	return result
}
