package smtp

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
	"time"
)

type SMTPClient struct {
	auth smtp.Auth
	host string
	port string
}

func New(username, password, host, port string) *SMTPClient {
	c := SMTPClient{
		host: host,
		port: port,
	}

	c.auth = smtp.PlainAuth("", username, password, host)

	return &c
}

func (c *SMTPClient) SendEmail(from string, to []string, subject, body string) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	sub := "Subject: " + subject + "!\n"
	msg := []byte(sub + mime + "\n" + body)

	if err := smtp.SendMail(c.getAddr(), c.auth, from, to, msg); err != nil {
		return err
	}

	return nil
}

func (c *SMTPClient) ParseTemplate(fileName string, data interface{}) (string, error) {
	d := &struct {
		Data interface{}
		Now  time.Time
	}{data, time.Now()}
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, d); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (c *SMTPClient) getAddr() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}
