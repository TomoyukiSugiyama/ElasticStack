#!/bin/sh

# ECR_URI=$(aws ecr describe-repositories --repository-names dev-repogitory --query 'repositories[].repositoryUri' --output text)
export AWS_VPC=$(aws ec2 describe-vpcs --filters 'Name=tag:f-iot.service.name,Values=vpc' --query 'Vpcs[0].VpcId' --output text)

# for elb_arn in $(aws elbv2 describe-load-balancers --query 'LoadBalancers[].LoadBalancerArn' --output text)
# do
#     tag=$(aws elbv2 describe-tags --resource-arns ${elb_arn} --query 'TagDescriptions[].[Tags[0].Value]' --output text)

#     if [ ${tag} = "nlb" ]; then
#         export ENDPOINT=$(aws elbv2 describe-load-balancers --load-balancer-arns ${elb_arn} --query 'LoadBalancers[0].DNSName' --output text)
#     fi
# done;

docker context use ecs

docker compose up