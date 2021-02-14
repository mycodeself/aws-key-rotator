package awsclient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type AWSIamClient interface {	
	ListAccessKeys(ctx context.Context, params *iam.ListAccessKeysInput, optFns ...func(*iam.Options)) (*iam.ListAccessKeysOutput, error)
	CreateAccessKey(ctx context.Context, params *iam.CreateAccessKeyInput, optFns ...func(*iam.Options)) (*iam.CreateAccessKeyOutput, error)
	UpdateAccessKey(ctx context.Context, params *iam.UpdateAccessKeyInput, optFns ...func(*iam.Options)) (*iam.UpdateAccessKeyOutput, error)
	DeleteAccessKey(ctx context.Context, params *iam.DeleteAccessKeyInput, optFns ...func(*iam.Options)) (*iam.DeleteAccessKeyOutput, error)
}

type Iam struct {
	client AWSIamClient
}

func CreateIamFromConfig(cfg aws.Config) *Iam {

	c := iam.NewFromConfig(cfg)

	awsIam := Iam{
		client: c,
	}

	return &awsIam
}

func (a *Iam) GetAccessKeysByUsername(ctx context.Context, username string) ([]types.AccessKeyMetadata, error) {
	input := &iam.ListAccessKeysInput{
		UserName: &username,
	}

	output, err := a.client.ListAccessKeys(ctx, input)

	if err != nil {
		return []types.AccessKeyMetadata{}, err
	}

	return output.AccessKeyMetadata, nil
}

func (a *Iam) CreateNewAccessKey(ctx context.Context, username string) (*types.AccessKey, error) {

	input := &iam.CreateAccessKeyInput{
		UserName: &username,
	}

	output, err := a.client.CreateAccessKey(ctx, input)

	if err != nil {
		return &types.AccessKey{}, err
	}

	return output.AccessKey, nil
}

func (a *Iam) DeactivateAccessKeyById(ctx context.Context, accessKeyId string, username string) error {
	input := &iam.UpdateAccessKeyInput{
		AccessKeyId: &accessKeyId,
		Status:      types.StatusTypeInactive,
		UserName:    &username,
	}

	_, err := a.client.UpdateAccessKey(ctx, input)

	if err != nil {
		return err
	}

	return nil
}

func (a *Iam) DeleteAccessKeyById(ctx context.Context, accessKeyId string, username string) error {
	input := &iam.DeleteAccessKeyInput{
		AccessKeyId: &accessKeyId,
		UserName:    &username,
	}

	_, err := a.client.DeleteAccessKey(ctx, input)

	if err != nil {
		return err
	}

	return nil
}
