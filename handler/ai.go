package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// GetAIConfig 获取用户AI配置
func GetAIConfig(c *gin.Context) {
	userID := c.GetUint("user_id")

	var aiConfig model.UserAIConfig
	if err := database.DB.Preload("Model").Where("user_id = ?", userID).First(&aiConfig).Error; err != nil {
		// 没有配置，返回空
		responseSuccess(c, gin.H{
			"user_id":     userID,
			"api_token":   "",
			"api_url":     "",
			"model_id":    0,
			"temperature": 0.7,
			"max_context": 10,
		})
		return
	}

	responseSuccess(c, aiConfig)
}

// UpdateAIConfig 更新用户AI配置
func UpdateAIConfig(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		APIToken    string  `json:"api_token"`
		APIUrl      string  `json:"api_url"`
		ModelID     uint    `json:"model_id"`
		ModelName   string  `json:"model_name"`
		Temperature float64 `json:"temperature"`
		MaxContext  int     `json:"max_context"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	var aiConfig model.UserAIConfig
	result := database.DB.Where("user_id = ?", userID).First(&aiConfig)
	if result.Error != nil {
		// 创建新配置
		aiConfig = model.UserAIConfig{
			UserID:      userID,
			APIToken:    req.APIToken,
			APIUrl:      req.APIUrl,
			ModelID:     req.ModelID,
			ModelName:   req.ModelName,
			Temperature: req.Temperature,
			MaxContext:  req.MaxContext,
		}
		if err := database.DB.Create(&aiConfig).Error; err != nil {
			responseErrorWithCode(c, http.StatusInternalServerError, "创建AI配置失败")
			return
		}
	} else {
		// 更新配置
		database.DB.Model(&aiConfig).Updates(map[string]interface{}{
			"api_token":   req.APIToken,
			"api_url":     req.APIUrl,
			"model_id":    req.ModelID,
			"model_name":  req.ModelName,
			"temperature": req.Temperature,
			"max_context": req.MaxContext,
		})
	}

	responseSuccess(c, nil)
}

// Chat AI对话（流式SSE输出、携带文章上下文、调用用户配置的AI模型）
func Chat(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		ArticleID uint   `json:"article_id"`
		Message   string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, "参数错误: "+err.Error())
		return
	}

	// 获取用户AI配置
	var aiConfig model.UserAIConfig
	if err := database.DB.Preload("Model").Where("user_id = ?", userID).First(&aiConfig).Error; err != nil {
		responseError(c, "请先配置AI参数")
		return
	}

	// 获取文章上下文
	var articleContext string
	if req.ArticleID > 0 {
		var article model.Article
		if err := database.DB.First(&article, req.ArticleID).Error; err == nil {
			articleContext = fmt.Sprintf("文章标题：%s\n文章内容：%s", article.Title, article.Content)
		}
	}

	// 保存用户消息
	userMsg := model.ChatMessage{
		UserID:    userID,
		ArticleID: req.ArticleID,
		Role:      "user",
		Content:   req.Message,
	}
	database.DB.Create(&userMsg)

	// 获取历史消息
	var history []model.ChatMessage
	database.DB.Where("user_id = ? AND article_id = ?", userID, req.ArticleID).
		Order("created_at DESC").Limit(aiConfig.MaxContext).Find(&history)

	// 构建请求消息
	messages := []map[string]string{}
	if articleContext != "" {
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": "你是一个智能助手，以下是相关文章内容供你参考：\n" + articleContext,
		})
	}
	for i := len(history) - 1; i >= 0; i-- {
		messages = append(messages, map[string]string{
			"role":    history[i].Role,
			"content": history[i].Content,
		})
	}
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": req.Message,
	})

	// 调用AI模型API
	apiURL := aiConfig.APIUrl
	if apiURL == "" && aiConfig.Model.ID > 0 {
		apiURL = aiConfig.Model.APIUrl
	}
	if apiURL == "" {
		responseError(c, "AI接口地址未配置")
		return
	}

	apiToken := aiConfig.APIToken
	modelName := aiConfig.ModelName
	if modelName == "" && aiConfig.Model.ID > 0 {
		modelName = aiConfig.Model.Name
	}

	// 构建请求体
	requestBody := map[string]interface{}{
		"model":       modelName,
		"messages":    messages,
		"stream":      true,
		"temperature": aiConfig.Temperature,
	}
	bodyBytes, _ := json.Marshal(requestBody)

	httpReq, _ := http.NewRequest("POST", apiURL, bytes.NewReader(bodyBytes))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiToken)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "AI接口调用失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// SSE流式输出
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	var fullContent string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if line[0] == ':' {
			continue
		}

		// 解析SSE数据
		if len(line) > 6 && line[:6] == "data: " {
			data := line[6:]
			if data == "[DONE]" {
				c.SSEvent("message", gin.H{"done": true})
				break
			}

			var sseData map[string]interface{}
			if err := json.Unmarshal([]byte(data), &sseData); err != nil {
				continue
			}

			if choices, ok := sseData["choices"].([]interface{}); ok && len(choices) > 0 {
				if choice, ok := choices[0].(map[string]interface{}); ok {
					if delta, ok := choice["delta"].(map[string]interface{}); ok {
						if content, ok := delta["content"].(string); ok {
							fullContent += content
							c.SSEvent("message", gin.H{"content": content})
							c.Writer.Flush()
						}
					}
				}
			}
		}
	}

	// 保存AI回复
	if fullContent != "" {
		assistantMsg := model.ChatMessage{
			UserID:    userID,
			ArticleID: req.ArticleID,
			Role:      "assistant",
			Content:   fullContent,
		}
		database.DB.Create(&assistantMsg)
	}
}

// GetAvailableModels 获取可用的AI模型列表
func GetAvailableModels(c *gin.Context) {
	var models []model.AIModel
	if err := database.DB.Where("enabled = ?", true).Order("id ASC").Find(&models).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取模型列表失败")
		return
	}
	responseSuccess(c, models)
}

// GetChatHistory 获取对话历史
func GetChatHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleID := c.Query("article_id")
	page, pageSize := getPaginationParams(c)

	query := database.DB.Model(&model.ChatMessage{}).Where("user_id = ?", userID)
	if articleID != "" {
		query = query.Where("article_id = ?", articleID)
	}

	var total int64
	query.Count(&total)

	var messages []model.ChatMessage
	offset := (page - 1) * pageSize
	if err := query.Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&messages).Error; err != nil {
		responseErrorWithCode(c, http.StatusInternalServerError, "获取对话历史失败")
		return
	}

	utils.ResponsePage(c, messages, total, page, pageSize)
}
