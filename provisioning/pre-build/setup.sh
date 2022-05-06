#!/bin/sh

# Common settings
SCRIPT_DIR=$(cd $(dirname $0); pwd)
WORK_DIR="$SCRIPT_DIR/../.."
alias cfn-stack-ops="$WORK_DIR/provisioning/helper-scripts/cfn-stack-ops.sh $1"

# KMS
# Create kms to encrypt/decrypt logs.
# cfn-stack-ops deploy cmk kms.yaml

# ECR
# Create elastic container registory.
# cfn-stack-ops deploy ecr $SCRIPT_DIR/ecr.yaml

# SSM
# Create parameters.
# Need to set github connection id and slack workspace/channel ids in secrets/*.yaml , before create parameters.
# $WORK_DIR/secrets/secret.sh

# S3
# Get s3 bucket name from ssm.
S3CfnBacketName=$(aws ssm get-parameter --name /dev/s3/cfn/BacketName --query "Parameter.Value" --output text)
S3LambdaBacketName=$(aws ssm get-parameter --name /dev/s3/lambda/BacketName --query "Parameter.Value" --output text)
# Create s3 backets





