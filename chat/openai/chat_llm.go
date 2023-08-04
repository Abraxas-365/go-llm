package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Abraxas-365/go-llm/chat"
)

type ChatOpenAI struct {
	Model       ChatModel
	Temperature uint32
	OpenAIKey   string
}

func NewChatOpenAI(options ...func(*ChatOpenAI)) chat.BaseChat {
	openaiKey, exists := os.LookupEnv("OPENAI_API_KEY")
	if !exists {
		openaiKey = ""
	}

	coa := &ChatOpenAI{
		Model:       Gpt3_5Turbo,
		Temperature: 0,
		OpenAIKey:   openaiKey,
	}

	for _, option := range options {
		option(coa)
	}

	return coa
}

func WithModel(model ChatModel) func(*ChatOpenAI) {
	return func(coa *ChatOpenAI) {
		coa.Model = model
	}
}

func WithTemperature(temperature uint32) func(*ChatOpenAI) {
	return func(coa *ChatOpenAI) {
		coa.Temperature = temperature
	}
}

func WithOpenAIKey(openaiKey string) func(*ChatOpenAI) {
	return func(coa *ChatOpenAI) {
		coa.OpenAIKey = openaiKey
	}
}

func (c *ChatOpenAI) Generate(messages [][]chat.BaseMessage) (*chat.AIMessage, error) {
	var flattenedMessages []Message
	for _, innerMessages := range messages {
		flattenedMessages = append(flattenedMessages, FromBaseMessages(innerMessages)...)
	}

	apiRequest := &Request{
		Model:    string(c.Model),
		Messages: flattenedMessages,
	}

	requestBody, err := json.Marshal(apiRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.OpenAIKey))

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return nil, NewOpenaiError(resp.StatusCode, bodyString)
	}

	var apiResponse Response
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	responseMessage := apiResponse.Choices[0]

	aiMessage := chat.NewAIMessage(responseMessage.Message.Content)

	return aiMessage.(*chat.AIMessage), nil
}

func (c *ChatOpenAI) Call(query string) (string, error) {

	apiRequest := &Request{
		Model:    string(c.Model),
		Messages: []Message{{Role: "user", Content: query}},
	}

	requestBody, err := json.Marshal(apiRequest)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.OpenAIKey))

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return "", NewOpenaiError(resp.StatusCode, bodyString)
	}

	var apiResponse Response
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return "", err
	}

	responseMessage := apiResponse.Choices[0]
	return responseMessage.Message.Content, nil
}
