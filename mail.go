package main

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

// SendEmail 发送邮件
func SendEmail(smtpHost, smtpPort, from, password string, to []string, subject, body string) error {
	// 邮件内容
	msg := "To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body

	// 如果是SSL/TLS端口（如465），使用TLS连接
	if smtpPort == "465" {
		// 创建TLS配置
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         smtpHost,
		}

		// 连接到SMTP服务器
		conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsConfig)
		if err != nil {
			return fmt.Errorf("连接SMTP服务器错误: %v", err)
		}

		// 创建SMTP客户端
		client, err := smtp.NewClient(conn, smtpHost)
		if err != nil {
			return fmt.Errorf("创建SMTP客户端错误: %v", err)
		}

		// 认证
		auth := smtp.PlainAuth("", from, password, smtpHost)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("认证错误: %v", err)
		}

		// 设置发件人
		if err := client.Mail(from); err != nil {
			return fmt.Errorf("错误设置发送方: %v", err)
		}

		// 设置收件人
		for _, addr := range to {
			if err := client.Rcpt(addr); err != nil {
				return fmt.Errorf("设置收件人错误: %v", err)
			}
		}

		// 发送邮件内容
		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("发送数据错误: %v", err)
		}
		_, err = w.Write([]byte(msg))
		if err != nil {
			return fmt.Errorf("写错误信息: %v", err)
		}
		err = w.Close()
		if err != nil {
			return fmt.Errorf("关闭数据写入器出错: %v", err)
		}

		// 关闭连接
		client.Quit()

		return nil
	}

	// 如果是普通端口（如587），使用smtp.SendMail
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
	if err != nil {
		return fmt.Errorf("发送邮件错误: %v", err)
	}

	return nil
}
