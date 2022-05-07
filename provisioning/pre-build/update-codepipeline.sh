#!/bin/sh

# Common settings
DEPLOY_ENV="dev" # set dev or prod

# Environment settings
SCRIPT_DIR=$(cd $(dirname $0); pwd)
WORK_DIR="${SCRIPT_DIR}/../.."
alias cfn-stack-ops="${WORK_DIR}/provisioning/helper-scripts/cfn-stack-ops.sh $1"

# CodeBuild, CodePipeline
# Get s3 bucket name from ssm.
S3CfnBucketName=$(aws ssm get-parameter --name /${DEPLOY_ENV}/s3/cfn/BucketName --query "Parameter.Value" --output text)

# Create codebuild and codepipeline.
cfn-stack-ops package ${WORK_DIR}/provisioning/pre-build/cfn/pre-build.yaml ${S3CfnBucketName} ${WORK_DIR}/provisioning/artifacts/pre-build-artifact.yaml
cfn-stack-ops deploy dev ${WORK_DIR}/provisioning/artifacts/pre-build-artifact.yaml