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

func NewProcess(config *config.RotationConfig, rotator *Rotator) *Process {
	p := &Process{
		config:    config,
		rotator:   rotator,
		notifiers: make(notifier.NotifiersMap),
	}

	// setup notifiers
	if util.StringInSlice("email", config.Notifiers) {
		p.notifiers["email"] = notifier.NewMailNotifier("templates/mail_result.txt")
	}

	if util.StringInSlice("slack", config.Notifiers) {
		p.notifiers["slack"] = notifier.NewSlackNotifierFromEnv()
	}

	return p
}

func (p *Process) Run(ctx context.Context) {

	var pResult notifier.ProcessResult

	for _, u := range p.config.AwsIamUsers {
		result := notifier.Result{
			Username: u.Username,
			ErrMsg:   "",
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
	for k, n := range p.notifiers {
		log.Info().Msgf("Sending result with %s notifier", k)
		err := n.NotifiyResult(ctx, result)
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}
}
