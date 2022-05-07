#!/bin/sh

docker context use default

aws ecr get-login-password --region ${AWS_REGION}| docker login --username AWS --password-stdin ${ECR_URI}

for SERVICE in logstash ecs-searchdomain-sidecar;
do
  docker image build -t ${ECR_URI}:${SERVICE} ${SERVICE}/
  docker image push ${ECR_URI}:${SERVICE}
done