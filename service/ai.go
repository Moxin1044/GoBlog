package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
)

type ChatRequest struct {
	Message   string `json:"message"`
	ArticleID uint   `json:"article_id"`
}

type AIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 构建AI对话上下文
func BuildChatContext(userID uint, req ChatRequest) ([]AIChatMessage, error) {
	messages := []AIChatMessage{}

	// 如果有文章ID，添加文章上下文
	if req.ArticleID > 0 {
		var article model.Article
		if err := database.DB.First(&article, req.ArticleID).Error; err == nil {
			messages = append(messages, AIChatMessage{
				Role:    "system",
				Content: fmt.Sprintf("以下是用户正在阅读的文章内容，请基于此内容回答问题：\n\n标题：%s\n\n%s", article.Title, article.Content),
			})
		}
	}

	// 加载历史对话
	var history []model.ChatMessage
	database.DB.Where("user_id = ? AND article_id = ?", userID, req.ArticleID).
		Order("created_at DESC").Limit(10).Find(&history)

	// 反转顺序
	for i := len(history) - 1; i >= 0; i-- {
		messages = append(messages, AIChatMessage{
			Role:    history[i].Role,
			Content: history[i].Content,
		})
	}

	// 添加当前消息
	messages = append(messages, AIChatMessage{
		Role:    "user",
		Content: req.Message,
	})

	return messages, nil
}

// 调用AI模型（流式）
func StreamChat(apiURL, apiToken string, model string, messages []AIChatMessage, temperature float64, maxTokens int) (<-chan string, error) {
	// 构建OpenAI兼容格式的请求体
	reqBody := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"stream":      true,
		"temperature": temperature,
	}
	if maxTokens > 0 {
		reqBody["max_tokens"] = maxTokens
	}

	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	ch := make(chan string, 100)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					ch <- fmt.Sprintf("data: {\"error\":\"%s\"}\n\n", err.Error())
				}
				return
			}
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				ch <- "data: [DONE]\n\n"
				return
			}
			var sseData map[string]interface{}
			if err := json.Unmarshal([]byte(data), &sseData); err != nil {
				continue
			}
			choices, ok := sseData["choices"].([]interface{})
			if !ok || len(choices) == 0 {
				continue
			}
			choice, ok := choices[0].(map[string]interface{})
			if !ok {
				continue
			}
			delta, ok := choice["delta"].(map[string]interface{})
			if !ok {
				continue
			}
			content, _ := delta["content"].(string)
			if content != "" {
				ch <- content
			}
		}
	}()

	return ch, nil
}

// 检查用户今日AI对话次数
func CheckDailyChatLimit(userID uint) (int, error) {
	var count int64
	err := database.DB.Model(&model.ChatMessage{}).
		Where("user_id = ? AND role = ? AND DATE(created_at) = CURDATE()", userID, "user").
		Count(&count).Error
	return int(count), err
}
