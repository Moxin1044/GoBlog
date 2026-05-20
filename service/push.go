package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/moxin/GoBlog/config"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
	"github.com/moxin/GoBlog/utils"
)

// 推送新文章通知给订阅用户
func PushNewArticle(article *model.Article, categoryIDs []uint) {
	var subs []model.Subscription
	// 查找订阅了该分类的用户
	database.DB.Joins("JOIN subscription_categories ON subscription_categories.subscription_id = subscriptions.id").
		Where("subscription_categories.category_id IN ?", categoryIDs).
		Find(&subs)

	for _, sub := range subs {
		go func(s model.Subscription) {
			if s.NotifyEmail && s.Email != "" {
				subject := fmt.Sprintf("新文章发布：%s", article.Title)
				body := fmt.Sprintf(`<h2>%s</h2><p>%s</p>`, article.Title, article.Summary)
				if err := utils.SendEmail(s.Email, subject, body); err != nil {
					log.Printf("Failed to send email to %s: %v", s.Email, err)
				}
			}
			if s.NotifyFeishu && s.FeishuToken != "" {
				sendFeishuMessage(s.FeishuToken, fmt.Sprintf("新文章发布：%s", article.Title))
			}
		}(sub)
	}
}

// 通知管理员
func NotifyAdmin(subject, content string) {
	// 飞书通知
	if config.AppConfig.Feishu.Enabled && config.AppConfig.Feishu.Token != "" {
		sendFeishuMessage(config.AppConfig.Feishu.Token, content)
	}
	// 邮件通知
	var admins []model.Admin
	database.DB.Where("status = ?", "active").Find(&admins)
	for _, admin := range admins {
		if admin.Email != "" {
			go func(email string) {
				if err := utils.SendEmail(email, subject, content); err != nil {
					log.Printf("Failed to notify admin %s: %v", email, err)
				}
			}(admin.Email)
		}
	}
}

// 发送飞书消息
func sendFeishuMessage(token, message string) error {
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v2/hook/%s", token)
	payload := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": message,
		},
	}
	body, _ := json.Marshal(payload)
	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bytes.NewReader(body)).
		Post(url)
	return err
}

// 服务器异常告警
func AlertServerAbnormal(metric, detail string) {
	subject := fmt.Sprintf("GoBlog服务器告警：%s异常", metric)
	content := fmt.Sprintf(`<h2>服务器告警</h2><p><strong>%s</strong>: %s</p>`, metric, detail)
	NotifyAdmin(subject, content)
}
