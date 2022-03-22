# ElasticStack
test
## Setting

```sh
# cat ~/.zshrc

alias cfn-stack-ops="/path/to/ElasticStack/provisioning/helper-scripts/cfn-stack-ops.sh $1"
alias s3-ops="/path/to/ElasticStack/provisioning/helper-scripts/s3-ops.sh $1"
```

## Deploy CFn

1. Create s3 bucket
```shell
$ s3-ops create elastic-stack-xxxxxxxxxxxxxxxxx-artifact
```

2. Package
```shell
$ mkdir provisioning/artifacts
$ cfn-stack-ops package provisioning/cfn/elastic-stack.yaml elastic-stack-xxxxxxxxxxxxxxxxx-artifact provisioning/artifacts/artifact.yaml
aws cloudformation package --template-file provisioning/cfn/elastic-stack.yaml --s3-bucket elastic-stack-xxxxxxxxxxxxxxxxx-artifact --output-template-file provisioning/artifacts/artifact.yaml
Uploading to XXXXXXXXXXXXXXXXXXXXXXXXXXXX.template  1669 / 1669.0  (100.00%)
Successfully packaged artifacts and wrote output template to file provisioning/artifacts/artifact.yaml.
Execute the following command to deploy the packaged template
aws cloudformation deploy --template-file /path/to/ElasticStack/provisioning/artifacts/artifact.yaml --stack-name <YOUR STACK NAME>
```

3. Deploy
```shell
$ cfn-stack-ops deploy test-stack provisioning/artifacts/artifact.yaml
aws cloudformation deploy --stack-name test-stack --template-file provisioning/artifacts/artifact.yaml

Waiting for changeset to be created..
Waiting for stack create/update to complete
Successfully created/updated stack - test-stack
```

## Deploy ECS

```sh
export AWS_REGION="ap-northeast-1"

# create a new repository and get repositoryUri
export ECR_URI=$(aws ecr create-repository \
  --repository-name test \
  --region $AWS_REGION \
  --query 'repository.repositoryUri' \
  --output text)

# get repositoryUri from existing repository
export ECR_URI=$(aws ecr describe-repositories \
  --repository-names test \
  --query 'repositories[].repositoryUri' \
  --output text)

export AWS_VPC=$(aws cloudformation describe-stacks \
  --stack-name test-stack \
  --query "Stacks[0].Outputs[?OutputKey=='VpcId'].OutputValue" \
  --output text)
```

deploy docker image
```bash
(cd provisioning/ecs ; ./deploy-container-image-to-ecr.sh)
```

check
```bash
$ aws ecr list-images --repository-name test | jq '.imageIds | .[].imageTag'
"ecs-serchdomain-sidecar"
"logstash"
```

deploy ecs
```bash
docker context use myecs
docker compose up
```

## Delete

* Delete s3 bucket
```shell
$ cfn-stack-ops delete elastic-stack-xxxxxxxxxxxxxxxxx-artifact
```

* Delete stack
```shell
$ cfn-stack-ops delete test-stack
```

## Utility

* S3
```shell
$ s3-ops 

Usage: /path/to/ElasticStack/provisioning/helper-scripts/s3-ops.sh MODE ARGS

Mode:     Args:
list      
create    s3-name
delete    s3-name
push      s3-name     file-path
```

* Cloudformation
```shell
$ cfn-stack-ops

Usage: /path/to/ElasticStack/provisioning/helper-scripts/cfn-stack-ops.sh MODE ARGS

Mode:     Args:
create    stack-name path-to-cfn-template-file [param1=val1 param2=val2]
update    stack-name path-to-cfn-template-file [param1=val1 param2=val2]
package   path-to-cfn-template-file s3-bucket output-template-file
deploy    stack-name path-to-cfn-template-filee
list      
describe  stack-name
validate  path-to-cfn-template-file
delete    stack-name
```

# SSH tunnel

```
ssh -i ~/.ssh/your-key.pem ec2-user@your-ec2-instance-public-ip -N -L 9200:vpc-domain-name.region.es.amazonaws.com:443
```