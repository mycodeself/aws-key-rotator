package target

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/mitchellh/mapstructure"
	"github.com/mycodeself/aws-key-rotator/pkg/circleci"
)

type CircleciProjectTarget struct {
	circleci *circleci.Client
}

type CircleciProjectTargetConfig struct {
	ProjectSlug            string `mapstructure:"project_slug"`
	AccessKeyIdVarName     string `mapstructure:"access_key_id_var_name"`
	SecretAccessKeyVarName string `mapstructure:"secret_access_key_var_name"`
}

func CreateCircleciProjectTarget(client *circleci.Client) *CircleciProjectTarget {
	return &CircleciProjectTarget{
		circleci: client,
	}
}

func (r *CircleciProjectTarget) Rotate(_ context.Context, key *types.AccessKey, config interface{}) error {
	c, err := r.parseConfig(config)
	if err != nil {
		return err
	}

	_, err = r.circleci.ProjectAddOrUpdateEnvVar(c.ProjectSlug, c.AccessKeyIdVarName, *key.AccessKeyId)
	if err != nil {
		// TODO: here critical error
		return err
	}

	_, err = r.circleci.ProjectAddOrUpdateEnvVar(c.ProjectSlug, c.SecretAccessKeyVarName, *key.SecretAccessKey)
	if err != nil {
		// TODO: ritical error, notify, necessary take action (key id updated but not access key)
		// how to handle this in a better way? implement rollback mechanism
		return err
	}

	return nil
}

func (r *CircleciProjectTarget) parseConfig(input interface{}) (*CircleciProjectTargetConfig, error) {
	var config CircleciProjectTargetConfig

	err := mapstructure.Decode(input, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
