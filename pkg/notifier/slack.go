package notifier

import (
	"context"
	"os"

	"github.com/slack-go/slack"
)

type SlackNotifier struct {
	client *slack.Client
	channelId string
}

func NewSlackNotifierFromEnv() *SlackNotifier {
	token := os.Getenv("SLACK_TOKEN")
	channel := os.Getenv("SLACK_CHANNEL")

	c := slack.New(token)
	
	s := SlackNotifier{
		client: c,
		channelId: channel,
	}

	return &s
}

func (n *SlackNotifier) NotifiyResult(ctx context.Context, result ProcessResult) error {
	_, _, err := n.client.PostMessage(
		n.channelId,
		slack.MsgOptionText("Some text", false),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		return err
	}

	return nil
}
