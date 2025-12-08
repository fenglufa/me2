package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/me2/ai/rpc/internal/config"
)

type Client struct {
	config     config.DeepseekConfig
	httpClient *http.Client
}

func NewClient(cfg config.DeepseekConfig) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature,omitempty"`
	MaxTokens   int32     `json:"max_tokens,omitempty"`
	TopP        float32   `json:"top_p,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int32 `json:"prompt_tokens"`
		CompletionTokens int32 `json:"completion_tokens"`
		TotalTokens      int32 `json:"total_tokens"`
	} `json:"usage"`
}

func (c *Client) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.ApiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回错误: %d, %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &chatResp, nil
}

type StreamChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

func (c *Client) ChatStream(ctx context.Context, req *ChatRequest, callback func(content string, done bool) error) error {
	req.Stream = true
	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.ApiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API返回错误: %d, %s", resp.StatusCode, string(body))
	}

	reader := io.Reader(resp.Body)
	buf := make([]byte, 4096)
	leftover := ""

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			data := leftover + string(buf[:n])
			lines := bytes.Split([]byte(data), []byte("\n"))

			for i, line := range lines {
				if i == len(lines)-1 {
					leftover = string(line)
					break
				}

				line = bytes.TrimSpace(line)
				if len(line) == 0 {
					continue
				}

				if bytes.HasPrefix(line, []byte("data: ")) {
					line = line[6:]
				}

				if string(line) == "[DONE]" {
					return callback("", true)
				}

				var chunk StreamChunk
				if err := json.Unmarshal(line, &chunk); err != nil {
					continue
				}

				if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
					if err := callback(chunk.Choices[0].Delta.Content, false); err != nil {
						return err
					}
				}
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取流失败: %w", err)
		}
	}

	return callback("", true)
}
