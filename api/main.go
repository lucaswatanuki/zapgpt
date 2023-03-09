package zapgpt

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Request struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

type Response struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Choices []Choice `json:"choices"`
}

func GenerateGPTText(query string) (string, error) {
	req := Request{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: query,
			},
		},
		MaxTokens: 150, // TODO env var
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqJson))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer token") //TODO set up token

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var resp Response
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
