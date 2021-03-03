package main

import (
	"context"
	"flag"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/mycodeself/aws-key-rotator/pkg/awsclient"
	"github.com/mycodeself/aws-key-rotator/pkg/circleci"
	"github.com/mycodeself/aws-key-rotator/pkg/config"
	"github.com/mycodeself/aws-key-rotator/pkg/rotator"
	"github.com/mycodeself/aws-key-rotator/pkg/target"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.TODO()

	if err := run(ctx); err != nil {
		log.Err(err).Msg(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	configPath := flag.String("config", "./config.yaml", "Path to the config file")
	flag.Parse()

	c, err := config.LoadFromYamlFile(*configPath)
	if err != nil {
		return errors.Wrap(err, "Error when loading configuration file")
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(""))
	if err != nil {
		return err
	}

	// initialize
	iamClient := awsclient.NewIamFromConfig(cfg)
	smClient := awsclient.NewSecretsManagerFromConfig(cfg)
	rot := rotator.NewRotator(iamClient)

	circleClient := circleci.NewClientFromEnv()
	if circleClient != nil {
		rot.AddRotationTarget("circleci_context", target.NewCircleciContextTarget(circleClient))
		rot.AddRotationTarget("circleci_project", target.NewCircleciProjectTarget(circleClient))
	}

	rot.AddRotationTarget("aws_secrets_manager_json", target.NewAwsSecretsManagerJsonTarget(smClient))

	p := rotator.NewProcess(c, rot)
	p.Run(ctx)

	return nil
}
