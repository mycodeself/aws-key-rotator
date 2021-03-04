package notifier

import (
	"context"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

type SlackNotifier struct {
	client    *slack.Client
	channelId string
}

func NewSlackNotifierFromEnv() *SlackNotifier {
	token := os.Getenv("SLACK_TOKEN")
	channel := os.Getenv("SLACK_CHANNEL")

	c := slack.New(token)

	s := SlackNotifier{
		client:    c,
		channelId: channel,
	}

	return &s
}

func (n *SlackNotifier) NotifiyResult(ctx context.Context, result ProcessResult) error {

	_, _, err := n.client.PostMessage(
		n.channelId,
		slack.MsgOptionText(n.createTextMessage(&result), false),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		return err
	}

	return nil
}

func (n *SlackNotifier) createTextMessage(result *ProcessResult) string {
	text := "*AWS Key Rotator results* \n-------------------\n"

	for _, r := range *result {
		line := fmt.Sprintf("User *%s*\n", r.Username)
		if r.ErrMsg != "" {
			line = line + fmt.Sprintf(">:exclamation::exclamation:Error: %s\n", r.ErrMsg)
		}

		if r.Rotated {
			line = line + ">:white_check_mark:Key rotation successfully completed!\n"
		} else {
			line = line + ">Rotation has not been performed\n"
		}

		text = text + line
	}

	return text
}
