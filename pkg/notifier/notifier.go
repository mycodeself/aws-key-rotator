package notifier

import (
	"context"
)

type Notifier interface {
	NotifiyResult(ctx context.Context, result ProcessResult) error
}

type NotifiersMap map[string]Notifier

type Result struct {
	Username string
	ErrMsg   string
	Rotated  bool
}

type ProcessResult []Result
