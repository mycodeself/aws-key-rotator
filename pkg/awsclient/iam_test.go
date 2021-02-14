package awsclient

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/mycodeself/aws-key-rotator/pkg/mock"
	"github.com/stretchr/testify/assert"
)


func TestGetAccessKeysByUsername(t *testing.T) {
	username := "user"

	c := &mock.AWSIamClientMock{
		DoListAccessKeys: func(ctx context.Context, params *iam.ListAccessKeysInput, optFns ...func(*iam.Options)) (*iam.ListAccessKeysOutput, error) {
			assert.Equal(t, username, *params.UserName)
			return &iam.ListAccessKeysOutput{}, nil
		},
	}
	i := &Iam{ client: c }
	
	_, err := i.GetAccessKeysByUsername(context.TODO(), username)	
	if err != nil {
		t.Error(err)
	}
}

func TestCreateNewAccessKey(t *testing.T) {
	username := "user"	
	c := &mock.AWSIamClientMock{
		DoCreateAccessKey: func(ctx context.Context, params *iam.CreateAccessKeyInput, optFns ...func(*iam.Options)) (*iam.CreateAccessKeyOutput, error) {
			assert.Equal(t, username, *params.UserName)
			return &iam.CreateAccessKeyOutput{}, nil
		},
	}
	i := &Iam{ client: c }
	
	_, err := i.CreateNewAccessKey(context.TODO(), username)
	if err != nil {
		t.Error(err)
	}
}

func TestDeactivateAccessKeyById(t *testing.T) {
	username := "user"
	accessKeyId := "ABCDE12345678"
	c := &mock.AWSIamClientMock{
		DoUpdateAccessKey: func(ctx context.Context, params *iam.UpdateAccessKeyInput, optFns ...func(*iam.Options)) (*iam.UpdateAccessKeyOutput, error) {
			assert.Equal(t, username, *params.UserName)
			assert.Equal(t, accessKeyId, *params.AccessKeyId)
			assert.Equal(t, types.StatusTypeInactive, params.Status)

			return &iam.UpdateAccessKeyOutput{}, nil
		},
	}
	i := &Iam{ client: c }

	err := i.DeactivateAccessKeyById(context.TODO(), accessKeyId, username)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteAccessKeyById(t *testing.T) {
	username := "user"
	accessKeyId := "ABCDE12345678"
	c := &mock.AWSIamClientMock{
		DoDeleteAccessKey: func(ctx context.Context, params *iam.DeleteAccessKeyInput, optFns ...func(*iam.Options)) (*iam.DeleteAccessKeyOutput, error) {
			assert.Equal(t, username, *params.UserName)
			assert.Equal(t, accessKeyId, *params.AccessKeyId)

			return &iam.DeleteAccessKeyOutput{}, nil
		},
	}
	i := &Iam{ client: c }

	err := i.DeleteAccessKeyById(context.TODO(), accessKeyId, username)
	if err != nil {
		t.Error(err)
	}
}