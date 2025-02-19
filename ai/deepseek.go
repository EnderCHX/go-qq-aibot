package ai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type DeepSeek struct {
	ApiUrl    string
	ApiKey    string
	Model     string
	SysPrompt string
}

type DeepSeekRequestBody struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream bool `json:"stream"`
}

type DeepSeekResponseBody struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs struct {
		} `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens          int `json:"prompt_tokens"`
		CompletionTokens      int `json:"completion_tokens"`
		TotalTokens           int `json:"total_tokens"`
		PromptCacheHitTokens  int `json:"prompt_cache_hit_tokens"`
		PromptCacheMissTokens int `json:"prompt_cache_miss_tokens"`
	} `json:"usage"`
}

func (d *DeepSeekRequestBody) ToJson() string {
	jsonBytes, err := json.Marshal(d)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

const (
	RoleSys       = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

func (d *DeepSeek) Init(apiurl, apikey, model, sysprompt string) {
	d.ApiUrl = apiurl
	d.ApiKey = apikey
	d.Model = model
	d.SysPrompt = sysprompt
}

func (d *DeepSeek) GetMessage(question string) (string, error) {
	fmt.Println("req api")
	requestBody := DeepSeekRequestBody{
		Model: d.Model,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    RoleSys,
				Content: d.SysPrompt,
			},
			{
				Role:    RoleUser,
				Content: question,
			},
		},
	}

	request, err := http.NewRequest("POST", d.ApiUrl, strings.NewReader(requestBody.ToJson()))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+d.ApiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		var responseBody DeepSeekResponseBody
		err := json.NewDecoder(response.Body).Decode(&responseBody)
		if err != nil {
			return "", err
		}
		// log.Println("message:", requestBody.Messages)
		// log.Println("responseBody:", responseBody.Choices)
		result := strings.ReplaceAll(responseBody.Choices[0].Message.Content, "\n", "")
		return result, nil
	}
	return "", nil
}
