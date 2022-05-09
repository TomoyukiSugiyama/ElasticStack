#!/bin/sh

for elb_arn in $(aws elbv2 describe-load-balancers --query 'LoadBalancers[].LoadBalancerArn' --output text)
do
    tag=$(aws elbv2 describe-tags --resource-arns ${elb_arn} --query 'TagDescriptions[].[Tags[0].Value]' --output text)

    if [ ${tag} = "nlb" ]; then
        export ENDPOINT=$(aws elbv2 describe-load-balancers --load-balancer-arns ${elb_arn} --query 'LoadBalancers[0].DNSName' --output text)
    fi
done;

echo ${ENDPOINT}
