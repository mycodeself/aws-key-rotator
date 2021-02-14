package mock

import (
	"context"

	awssm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AWSSecretsManagerClientMock struct {
	DoGetSecretValue func(ctx context.Context, params *awssm.GetSecretValueInput, optFns ...func(*awssm.Options)) (*awssm.GetSecretValueOutput, error)
	DoUpdateSecret func(ctx context.Context, params *awssm.UpdateSecretInput, optFns ...func(*awssm.Options)) (*awssm.UpdateSecretOutput, error)
}

func (m *AWSSecretsManagerClientMock) GetSecretValue(ctx context.Context, params *awssm.GetSecretValueInput, optFns ...func(*awssm.Options)) (*awssm.GetSecretValueOutput, error) {
	return m.DoGetSecretValue(ctx, params, optFns...)
}

func (m *AWSSecretsManagerClientMock) UpdateSecret(ctx context.Context, params *awssm.UpdateSecretInput, optFns ...func(*awssm.Options)) (*awssm.UpdateSecretOutput, error) {
	return m.DoUpdateSecret(ctx, params, optFns...)
}