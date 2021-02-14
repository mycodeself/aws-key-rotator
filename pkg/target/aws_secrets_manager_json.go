package target

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/mitchellh/mapstructure"
	"github.com/mycodeself/aws-key-rotator/pkg/awsclient"
)

type AwsSecretsManagerJsonTarget struct {
	sm *awsclient.SecretsManager
}

type AwsSecretsManagerJsonConfig struct {
	SecretArn               string `mapstructure:"secret_arn"`
	KmsKeyId                string `mapstructure:"kms_key_id"`
	AccessKeyIdProperty     string `mapstructure:"access_key_id_property"`
	SecretAccessKeyProperty string `mapstructure:"secret_access_key_property"`
}

func CreateAwsSecretsManagerJsonTarget(sm *awsclient.SecretsManager) *AwsSecretsManagerJsonTarget {
	a := AwsSecretsManagerJsonTarget{
		sm: sm,
	}

	return &a
}

func (r *AwsSecretsManagerJsonTarget) Rotate(ctx context.Context, key *types.AccessKey, config interface{}) error {
	c, err := r.parseConfig(config)
	if err != nil {
		return err
	}

	var secretVal map[string]interface{}

	val, err := r.sm.GetSecretValue(ctx, c.SecretArn)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), &secretVal)
	if err != nil {
		return err
	}

	// update secret properties
	secretVal[c.AccessKeyIdProperty] = key.AccessKeyId
	secretVal[c.SecretAccessKeyProperty] = key.SecretAccessKey

	secretStr, err := json.Marshal(secretVal)
	if err != nil {
		return err
	}

	// save secret string
	err = r.sm.UpdateSecretValue(ctx, c.SecretArn, string(secretStr), c.KmsKeyId)
	if err != nil {
		return err
	}

	return nil
}

func (r *AwsSecretsManagerJsonTarget) parseConfig(input interface{}) (*AwsSecretsManagerJsonConfig, error) {
	var config AwsSecretsManagerJsonConfig

	err := mapstructure.Decode(input, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
