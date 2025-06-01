package Aichat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

const (
	apiURL    = "https://api.deepseek.com/v1/chat/completions" // 请根据实际API地址修改
	apiKey    = "sk-fe7e8221cd8a4c2fb5d029d03cc2bc9e"          // 替换为你的API密钥
	modelName = "deepseek-chat"                                // 替换为实际模型名称
)

// 定义请求结构体
type ChatRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 定义响应结构体
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func AiChat(c *gin.Context) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	var msg Message
	err := c.ShouldBind(&msg)
	if err != nil {
		return
	}
	// 构造请求
	reqBody := ChatRequest{
		Model: modelName,
		Messages: []Message{
			{
				Role:    "user",
				Content: msg.Content,
			},
		},
		MaxTokens: 500,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("Error marshaling request:", err)
		return
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(500, gin.H{
			"Error reading response": err,
			"code":                   500,
		})
		return
	}

	// 解析响应
	var chatResp ChatResponse
	if err = json.Unmarshal(body, &chatResp); err != nil {
		c.JSON(500, gin.H{
			"Error unmarshaling response": err,
			"code":                        500,
		})
		return
	}

	// 处理响应
	if len(chatResp.Choices) > 0 {
		c.JSON(200, gin.H{
			"code": 200,
			"data": chatResp.Choices[0].Message.Content})
	} else if chatResp.Error.Message != "" {
		c.JSON(200, gin.H{
			"code": 200,
			"data": "API返回错误：" + chatResp.Error.Message,
		})
	} else {
		c.JSON(200, gin.H{
			"code": 200,
			"data": "未知响应：" + string(body),
		})
	}
}
