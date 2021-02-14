# under development...

# AWS Key Rotator

AWS Key Rotator is a golang program that helps you with the repetitive task of rotating the credentials of your AWS IAM system accounts.

The process is simple, configure the credentials that you want to rotate and those targets that should be rotated when rotating these credentials.

A target is a entity that makes use of these credentials and must be updated (often manually) to continue working properly. As an example, an AWS IAM user used in CircleCI to upload an image to the ECR. When the AWS credentials are rotated, the environment variables in CircleCI must be updated with the new ones.

# Usage

## Available targets

#### AWS Secrets Manager

AWS Secrets Manager targets automatically update a secret stored in AWS Secrets Manager service.

_Note: This uses the default AWS credentials configured in the system, so no extra configuration is needed_

##### AWS Secrets Manager JSON Target

This target will automatically update a secret stored in AWS Secrets Manager in JSON format. It updates the Access Key Id and Secret Access Key in the specified JSON properties and keeps the rest of the JSON.

```yaml
aws_iam_users:
  - username: user-to-rotate
    days: 60
    targets:
      - aws_secrets_manager_json:
          secret_arn: arn:aws:secretsmanager:eu-west-1:123456789:secret:mysecret-12345
          access_key_id_property: AWS_SECRET_KEY_ID
          secret_access_key_property: AWS_SECRET_ACCESS_KEY
          kms_key_id: # (Optional) Specifies an updated ARN or alias of the AWS KMS customer master key
```

#### CircleCI

Ensure `CIRCLECI_TOKEN` environment variable is present with a valid API token to access the projects or contexts you want to automatically update, see https://circleci.com/docs/2.0/managing-api-tokens/

##### CircleCI Context Target

```yaml
aws_iam_users:
  - username: user-to-rotate
    days: 60
    targets:
      - circleci_context:
          context_id: 8cea5754-907d-4425-9b7f-8493de1efbfa
          access_key_id_var_name: AWS_SECRET_KEY_ID
          secret_access_key_var_name: AWS_SECRET_ACCESS_KEY
```

##### CircleCI Project Target

```yaml
aws_iam_users:
  - username: user-to-rotate
    days: 60
    targets:
      - circleci_project:
          project_slug: github/user/project
          access_key_id_var_name: AWS_SECRET_KEY_ID
          secret_access_key_var_name: AWS_SECRET_ACCESS_KEY
```
