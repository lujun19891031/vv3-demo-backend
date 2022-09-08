package utils

import (
	"math/rand"
	"net/smtp"
	"strings"
	"time"
)

// smtp服务器配置
var (
	SmtpHost     = "smtp.263.net"
	SmtpPort     = "25"
	SmtpUser     = "system@penaviconb.com"
	SmtpPassword = "net263"
)

// 邮件发送
func SendMail(address []string, subject, body string) (err error) {
	// 发送邮件
	auth := smtp.PlainAuth("", SmtpUser, SmtpPassword, SmtpHost)
	to := address
	msg := []byte("To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + SmtpUser + "<" + SmtpUser + ">\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)
	err = smtp.SendMail(SmtpHost+":"+SmtpPort, auth, SmtpUser, to, msg)
	return
}

// 随机生成验证码
func GetRandomCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// 随机生成字符串
	var bytes = make([]byte, length)
	for i := 0; i < length; i++ {
		b := rand.Intn(len(str))
		bytes[i] = str[b]
	}
	return string(bytes)
}
