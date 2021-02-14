package awsclient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AWSSecretsManagerClient interface {
	GetSecretValue(ctx context.Context, params *awssm.GetSecretValueInput, optFns ...func(*awssm.Options)) (*awssm.GetSecretValueOutput, error)
	UpdateSecret(ctx context.Context, params *awssm.UpdateSecretInput, optFns ...func(*awssm.Options)) (*awssm.UpdateSecretOutput, error)
}

type SecretsManager struct {
	client AWSSecretsManagerClient
}

func CreateSecretsManager(cfg aws.Config) *SecretsManager {

	c := awssm.NewFromConfig(cfg)

	sm := &SecretsManager{
		client: c,
	}

	return sm
}

func (sm *SecretsManager) GetSecretValue(ctx context.Context, secretId string) (string, error) {
	input := &awssm.GetSecretValueInput{
		SecretId: &secretId,
	}

	val, err := sm.client.GetSecretValue(ctx, input)

	if err != nil {
		return "", err
	}

	return *val.SecretString, nil
}

func (sm *SecretsManager) UpdateSecretValue(ctx context.Context, secretId, value, kmsKeyId string) error {
	input := &awssm.UpdateSecretInput{
		SecretId:     &secretId,
		KmsKeyId:     &kmsKeyId,
		SecretString: &value,
	}

	_, err := sm.client.UpdateSecret(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
