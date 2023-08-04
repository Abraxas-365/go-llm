package chat

type BaseChat interface {
	Generate(messages []BaseMessage)
	Call(query string) (string, error)
}
