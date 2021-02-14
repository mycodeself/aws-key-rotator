package rotator

import (
	"context"

	"github.com/mycodeself/aws-key-rotator/pkg/config"
	"github.com/mycodeself/aws-key-rotator/pkg/notifier"
	"github.com/mycodeself/aws-key-rotator/pkg/util"

	"github.com/rs/zerolog/log"
)

type Process struct {
	config    *config.RotationConfig
	rotator   *Rotator
	notifiers notifier.NotifiersMap
}

func CreateProcess(config *config.RotationConfig, rotator *Rotator) *Process {
	p := &Process{
		config:    config,
		rotator:   rotator,
		notifiers: make(notifier.NotifiersMap),
	}

	// setup notifiers
	if util.StringInSlice("email", config.Notifiers) {
		p.notifiers["email"] = notifier.CreateMailNotifier("templates/mail_result.txt")
	}

	return p
}

func (p *Process) Run(ctx context.Context) {

	var pResult notifier.ProcessResult

	for _, u := range p.config.AwsIamUsers {
		result := notifier.Result{
			Username: u.Username,
		}

		r, err := p.rotator.RotateAwsUser(ctx, &u, p.config.SafeMode)
		if err != nil {
			log.Error().Msg(err.Error())
			result.ErrMsg = err.Error()
		}

		result.Rotated = r

		pResult = append(pResult, result)
	}

	p.dispatchNotifiers(ctx, pResult)

}

func (p *Process) dispatchNotifiers(ctx context.Context, result notifier.ProcessResult) {
	for _, n := range p.notifiers {
		err := n.NotifiyResult(ctx, result)
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}
}
