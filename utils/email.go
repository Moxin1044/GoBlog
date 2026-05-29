package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/moxin/GoBlog/config"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
)

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

func getSMTPConfig() *SMTPConfig {
	var configs []model.SiteConfig
	database.DB.Where("`key` IN ?", []string{"smtp_host", "smtp_port", "smtp_user", "smtp_password", "smtp_sender"}).Find(&configs)

	kv := make(map[string]string)
	for _, c := range configs {
		kv[c.Key] = c.Value
	}

	if kv["smtp_host"] != "" {
		port, _ := strconv.Atoi(kv["smtp_port"])
		if port == 0 {
			port = 465
		}
		from := kv["smtp_sender"]
		if from == "" {
			from = kv["smtp_user"]
		}
		return &SMTPConfig{
			Host:     kv["smtp_host"],
			Port:     port,
			User:     kv["smtp_user"],
			Password: kv["smtp_password"],
			From:     from,
		}
	}

	cfg := config.AppConfig.SMTP
	if cfg.Host == "" {
		return nil
	}
	from := cfg.From
	if from == "" {
		from = cfg.User
	}
	return &SMTPConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		From:     from,
	}
}

func SendEmail(to, subject, body string) error {
	cfg := getSMTPConfig()
	if cfg == nil {
		return fmt.Errorf("SMTP not configured")
	}

	msg := strings.Join([]string{
		"From: " + cfg.From,
		"To: " + to,
		"Subject: " + subject,
		"Content-Type: text/html; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)

	tlsConfig := &tls.Config{
		ServerName: cfg.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return err
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(cfg.From); err != nil {
		return err
	}

	if err = client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}

func SendVerificationCode(to, code string) error {
	subject := "GoBlog 验证码"
	body := fmt.Sprintf(`<h2>您的验证码是：%s</h2><p>验证码有效期为5分钟，请尽快使用。</p>`, code)
	return SendEmail(to, subject, body)
}
