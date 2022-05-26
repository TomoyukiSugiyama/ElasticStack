#!/bin/sh

# Environment settings
SCRIPT_DIR=$(cd $(dirname $0); pwd)
WORK_DIR="${SCRIPT_DIR}/../.."
alias cfn-stack-ops="${WORK_DIR}/provisioning/helper-scripts/cfn-stack-ops.sh $1"

export ECR_URI=$(aws ecr describe-repositories --repository-names dev-repogitory --query 'repositories[].repositoryUri' --output text)

${SCRIPT_DIR}/ecs/deploy-container-image-to-ecr.sh

AWS_VPC=$(aws ec2 describe-vpcs --filters 'Name=tag:f-iot.service.name,Values=vpc' --query 'Vpcs[0].VpcId' --output text)
public_subnet_1=$(aws ec2 describe-subnets --filters "Name=tag:f-iot.service.name,Values=public-subnet-1" --query 'Subnets[0].SubnetId' --output text)
public_subnet_2=$(aws ec2 describe-subnets --filters "Name=tag:f-iot.service.name,Values=public-subnet-2" --query 'Subnets[0].SubnetId' --output text)
SubnetIds="${public_subnet_1},${public_subnet_2}"

cfn-stack-ops deploy filebeat ${SCRIPT_DIR}/cfn/filebeat.yaml VpcId=${AWS_VPC} SubnetIds=${SubnetIds}
