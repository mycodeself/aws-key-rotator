package notifier

import (
	"context"
	"os"
	"strings"

	"github.com/mycodeself/aws-key-rotator/pkg/smtp"
)

type MailNotifier struct {
	smtp         smtp.SMTPClient
	templateFile string
	to           []string
	from         string
}

func CreateMailNotifier(templateFile string) *MailNotifier {
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	to := strings.Split(os.Getenv("SMTP_TO"), ",")
	from := os.Getenv("SMTP_FROM")
	port := os.Getenv("SMTP_PORT")

	s := smtp.Create(username, password, host, port)
	n := MailNotifier{
		smtp:         *s,
		templateFile: templateFile,
		to:           to,
		from:         from,
	}

	return &n
}

func (n *MailNotifier) NotifiyResult(ctx context.Context, result ProcessResult) error {
	content, err := n.smtp.ParseTemplate(n.templateFile, result)
	subject := "AWS Key Rotator Friend - Rotation result"

	if err != nil {
		return err
	}

	return n.smtp.SendEmail(n.from, n.to, subject, content)
}
