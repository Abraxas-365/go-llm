package embedding

type BaseEmbedder interface {
	EmbedDocuments(documents []string) ([][]float64, error)
	EmbedQuery(text string) ([]float64, error)
}
