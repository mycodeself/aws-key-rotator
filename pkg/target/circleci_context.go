package target

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/mitchellh/mapstructure"
	"github.com/mycodeself/aws-key-rotator/pkg/circleci"
)

type CircleciContextTarget struct {
	circleci *circleci.Client
}

type CircleciContextTargetConfig struct {
	ContextId              string `mapstructure:"context_id"`
	AccessKeyIdVarName     string `mapstructure:"access_key_id_var_name"`
	SecretAccessKeyVarName string `mapstructure:"secret_access_key_var_name"`
}

func CreateCircleciContextTarget(client *circleci.Client) *CircleciContextTarget {
	c := CircleciContextTarget{
		circleci: client,
	}

	return &c
}

func (r *CircleciContextTarget) Rotate(_ context.Context, key *types.AccessKey, config interface{}) error {
	c, err := r.parseConfig(config)
	if err != nil {
		return err
	}

	_, err = r.circleci.ContextAddOrUpdateEnvVar(c.ContextId, c.AccessKeyIdVarName, *key.AccessKeyId)
	if err != nil {
		// TODO: here critical error
		return err
	}

	_, err = r.circleci.ContextAddOrUpdateEnvVar(c.ContextId, c.SecretAccessKeyVarName, *key.SecretAccessKey)
	if err != nil {
		// TODO: critical error, notify, necessary take action (key id updated but not access key)
		// how to handle this in a better way? implement rollback mechanism
		return err
	}

	return nil
}

func (r *CircleciContextTarget) parseConfig(input interface{}) (*CircleciContextTargetConfig, error) {
	var config CircleciContextTargetConfig

	err := mapstructure.Decode(input, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
