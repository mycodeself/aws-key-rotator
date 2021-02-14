package awsclient

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/mycodeself/aws-key-rotator/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetSecretValue(t *testing.T) {
	secretId := "arn:aws:secretsmanager:eu-west-1:123456789:secret:mysecret-12345"
	c := &mock.AWSSecretsManagerClientMock{
		DoGetSecretValue: func(ctx context.Context, params *awssm.GetSecretValueInput, optFns ...func(*awssm.Options)) (*awssm.GetSecretValueOutput, error) {
			assert.Equal(t, secretId, *params.SecretId)

			return &awssm.GetSecretValueOutput{ SecretString: aws.String("secret")}, nil
		},
	}
	sm := &SecretsManager{ client: c } 

	_, err := sm.GetSecretValue(context.TODO(), secretId)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateSecretValue(t *testing.T) {
	secretId := "arn:aws:secretsmanager:eu-west-1:123456789:secret:mysecret-12345"
	value := "secret"

	c := &mock.AWSSecretsManagerClientMock{
		DoUpdateSecret: func(ctx context.Context, params *awssm.UpdateSecretInput, optFns ...func(*awssm.Options)) (*awssm.UpdateSecretOutput, error) {
			assert.Equal(t, secretId, *params.SecretId)
			assert.Equal(t, value, *params.SecretString)

			return &awssm.UpdateSecretOutput{}, nil
		},
	}
	sm := &SecretsManager{ client: c } 

	err := sm.UpdateSecretValue(context.TODO(), secretId, value, "")
	if err != nil {
		t.Error(err)
	}
}