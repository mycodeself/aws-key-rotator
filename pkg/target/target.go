package target

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type RotationTarget interface {
	Rotate(ctx context.Context, newKey *types.AccessKey, config interface{}) error
}

type RotationTargetMap map[string]RotationTarget
