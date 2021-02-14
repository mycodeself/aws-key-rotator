package rotator

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/mycodeself/aws-key-rotator/pkg/awsclient"
	"github.com/mycodeself/aws-key-rotator/pkg/config"
	"github.com/mycodeself/aws-key-rotator/pkg/target"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Rotator struct {
	iam             *awsclient.Iam
	rotationTargets target.RotationTargetMap
}

func Create(iam *awsclient.Iam) *Rotator {
	r := Rotator{
		iam:             iam,
		rotationTargets: make(target.RotationTargetMap),
	}

	return &r
}

func (r *Rotator) AddRotationTarget(name string, t target.RotationTarget) {
	r.rotationTargets[name] = t
}

// Performs the rotation of AWS IAM access key
func (r *Rotator) RotateAwsUser(ctx context.Context, u *config.AwsIamUser, safeMode bool) (bool, error) {
	keys, err := r.iam.GetAccessKeysByUsername(ctx, u.Username)

	if err != nil {
		return false, errors.Wrap(err, "Error when retrieving access keys by username from aws")
	}

	if len(keys) != 1 {
		return false, errors.New(fmt.Sprintf("Retrieved %d keys from aws, expected 1 key", len(keys)))
	}

	key := &keys[0]

	if !r.shouldRotateKey(key, u.Days) {
		log.Info().Msgf("Username %s doesn't need to be rotated", u.Username)
		return false, nil
	}

	log.Info().Msgf("Rotating key %s", *key.AccessKeyId)

	newKey, err := r.iam.CreateNewAccessKey(ctx, u.Username)

	if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("Error when creating new access key for user %s", u.Username))
	}

	log.Info().Msgf("New access key has been created")

	// rotate targets, circleci, aws secrets manager...
	for _, t := range u.Targets {
		for k, c := range t {
			if target, ok := r.rotationTargets[k]; ok {
				err = target.Rotate(ctx, newKey, c)
				if err != nil {
					return false, errors.Wrap(err, fmt.Sprintf("Error rotating with %s target", k))
				}

			} else {
				log.Warn().Msgf("Target type %s doesn't exists, review your configuration", k)
			}
		}
	}

	err = r.postRotation(ctx, key, safeMode)

	if err != nil {
		return false, errors.Wrap(err, "An error occurred during post rotation")
	}

	return true, nil
}

func (r *Rotator) shouldRotateKey(key *types.AccessKeyMetadata, days int) bool {
	diff := time.Since(*key.CreateDate)
	daysDiff := int(math.Floor(diff.Hours() / 24))

	return daysDiff >= days
}

func (r *Rotator) postRotation(ctx context.Context, oldKey *types.AccessKeyMetadata, safeMode bool) error {
	log.Info().Msgf("Deactivating previous access key %s", *oldKey.AccessKeyId)

	err := r.iam.DeactivateAccessKeyById(ctx, *oldKey.AccessKeyId, *oldKey.UserName)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error when deactivating access key %s", *oldKey.AccessKeyId))
	}

	if !safeMode {
		log.Info().Msgf("Deleting previous access key %s", *oldKey.AccessKeyId)

		err = r.iam.DeleteAccessKeyById(ctx, *oldKey.AccessKeyId, *oldKey.UserName)

		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Error when deleting access key %s", *oldKey.AccessKeyId))
		}
	}

	return nil
}
