package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/moxin/GoBlog/config"
)

func SendEmail(to, subject, body string) error {
	cfg := config.AppConfig.SMTP
	if cfg.Host == "" {
		return fmt.Errorf("SMTP not configured")
	}

	from := cfg.From
	if from == "" {
		from = cfg.User
	}

	msg := strings.Join([]string{
		"From: " + from,
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

	if err = client.Mail(from); err != nil {
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
