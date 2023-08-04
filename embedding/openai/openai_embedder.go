package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Abraxas-365/go-llm/embedding"
)

type OpenAiEmbedder struct {
	Model     string
	OpenAIKey string
}

type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingsResponse struct {
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

func NewOpenAiEmbedder(options ...func(*OpenAiEmbedder)) embedding.BaseEmbedder {
	openaiKey, exists := os.LookupEnv("OPENAI_API_KEY")
	if !exists {
		openaiKey = ""
	}

	oae := &OpenAiEmbedder{
		Model:     "text-embedding-ada-002",
		OpenAIKey: openaiKey,
	}

	for _, option := range options {
		option(oae)
	}

	return oae
}

func WithModel(model string) func(*OpenAiEmbedder) {
	return func(oae *OpenAiEmbedder) {
		oae.Model = model
	}
}

func (o *OpenAiEmbedder) EmbedDocuments(documents []string) ([][]float64, error) {
	payload := map[string]interface{}{
		"input": documents,
		"model": o.Model,
	}

	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.OpenAIKey)

	client := &http.Client{Timeout: time.Minute * 3}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("non-200 response from OpenAI API")
	}

	body, _ := ioutil.ReadAll(res.Body)

	var data EmbeddingsResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	var embeddings [][]float64
	for _, embedding := range data.Data {
		embeddings = append(embeddings, embedding.Embedding)
	}

	return embeddings, nil
}

func (o *OpenAiEmbedder) EmbedQuery(text string) ([]float64, error) {
	payload := map[string]interface{}{
		"input": text,
		"model": o.Model,
	}

	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.OpenAIKey)

	client := &http.Client{Timeout: time.Minute * 3}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("non-200 response from OpenAI API")
	}

	body, _ := ioutil.ReadAll(res.Body)

	var data EmbeddingsResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if len(data.Data) == 0 {
		return nil, errors.New("no embeddings in response")
	}

	return data.Data[0].Embedding, nil
}
