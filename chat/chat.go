package chat

type BaseChat interface {
	Generate(messages [][]BaseMessage) (*AIMessage, error)
	Call(query string) (string, error)
}
