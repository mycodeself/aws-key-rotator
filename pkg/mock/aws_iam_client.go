package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type AWSIamClientMock struct {
	DoListAccessKeys func(ctx context.Context, params *iam.ListAccessKeysInput, optFns ...func(*iam.Options)) (*iam.ListAccessKeysOutput, error)
	DoCreateAccessKey func(ctx context.Context, params *iam.CreateAccessKeyInput, optFns ...func(*iam.Options)) (*iam.CreateAccessKeyOutput, error)
	DoUpdateAccessKey func(ctx context.Context, params *iam.UpdateAccessKeyInput, optFns ...func(*iam.Options)) (*iam.UpdateAccessKeyOutput, error)
	DoDeleteAccessKey func(ctx context.Context, params *iam.DeleteAccessKeyInput, optFns ...func(*iam.Options)) (*iam.DeleteAccessKeyOutput, error)
}

func (m *AWSIamClientMock) ListAccessKeys(ctx context.Context, params *iam.ListAccessKeysInput, optFns ...func(*iam.Options)) (*iam.ListAccessKeysOutput, error) {
	return m.DoListAccessKeys(ctx,params, optFns...)
}

func (m *AWSIamClientMock) CreateAccessKey(ctx context.Context, params *iam.CreateAccessKeyInput, optFns ...func(*iam.Options)) (*iam.CreateAccessKeyOutput, error) {
	return m.DoCreateAccessKey(ctx, params, optFns...)
}

func (m *AWSIamClientMock) UpdateAccessKey(ctx context.Context, params *iam.UpdateAccessKeyInput, optFns ...func(*iam.Options)) (*iam.UpdateAccessKeyOutput, error) {
	return m.DoUpdateAccessKey(ctx, params, optFns...)
}

func (m *AWSIamClientMock) DeleteAccessKey(ctx context.Context, params *iam.DeleteAccessKeyInput, optFns ...func(*iam.Options)) (*iam.DeleteAccessKeyOutput, error) {
	return m.DoDeleteAccessKey(ctx, params, optFns...)
}
