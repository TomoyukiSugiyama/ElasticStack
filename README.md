# ElasticStack

## Deployment
![](https://codebuild.ap-northeast-1.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiRnc1MjZTTS9oMGR6QXpCdVlKbVhyR0xROHlab2E1MUE5VHdJVTJ5aXI3RVJyWnJzUEVqVkpkeHNiUkdSUFNxTGJlaExpY3BDVmMxYjA4RkltQ0pIMWVvPSIsIml2UGFyYW1ldGVyU3BlYyI6IlpkVTZMMVVjZmhwMHF3MjkiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

## Setting
```bash
# cat ~/.zshrc

alias cfn-stack-ops="/path/to/ElasticStack/provisioning/helper-scripts/cfn-stack-ops.sh $1"
alias s3-ops="/path/to/ElasticStack/provisioning/helper-scripts/s3-ops.sh $1"
```

## Build CFn
1. Create s3 bucket
```bash
$ s3-ops create elastic-stack-xxxxxxxxxxxxxxxxx-artifact
$ export S3_BUCKET_CFN_NAME="elastic-stack-xxxxxxxxxxxxxxxxx-artifact"
```

2. Package
```bash
$ mkdir provisioning/artifacts
$ cfn-stack-ops package provisioning/cfn/elastic-stack.yaml $S3_BUCKET_CFN_NAME provisioning/artifacts/artifact.yaml
aws cloudformation package --template-file provisioning/cfn/elastic-stack.yaml --s3-bucket elastic-stack-xxxxxxxxxxxxxxxxx-artifact --output-template-file provisioning/artifacts/artifact.yaml
Uploading to XXXXXXXXXXXXXXXXXXXXXXXXXXXX.template  1669 / 1669.0  (100.00%)
Successfully packaged artifacts and wrote output template to file provisioning/artifacts/artifact.yaml.
Execute the following command to deploy the packaged template
aws cloudformation deploy --template-file /path/to/ElasticStack/provisioning/artifacts/artifact.yaml --stack-name <YOUR STACK NAME>
```

3. Push cfn to s3 bucket
```bash
$ cfn-stack-ops deploy test-stack provisioning/artifacts/artifact.yaml
aws cloudformation deploy --stack-name test-stack --template-file provisioning/artifacts/artifact.yaml

Waiting for changeset to be created..
Waiting for stack create/update to complete
Successfully created/updated stack - test-stack
```

Create stack from s3 bucket
```bash
$ cfn-stack-ops create test-stack $S3_BUCKET_CFN_NAME
```

Update stack from s3 bucket
```bash
$ cfn-stack-ops update test-stack $S3_BUCKET_CFN_NAME
```

## Build ECS
1. Create s3 bucket and set env
```bash
export AWS_REGION="ap-northeast-1"

# create a new repository and get repositoryUri
export ECR_URI=$(aws ecr create-repository \
  --repository-name f-iot-rep \
  --region $AWS_REGION \
  --query 'repository.repositoryUri' \
  --output text)

# get repositoryUri from existing repository
export ECR_URI=$(aws ecr describe-repositories \
  --repository-names f-iot-rep \
  --query 'repositories[].repositoryUri' \
  --output text)

export AWS_VPC=$(aws cloudformation describe-stacks \
  --stack-name test-stack \
  --query "Stacks[0].Outputs[?OutputKey=='VpcId'].OutputValue" \
  --output text)
```

2. Build images and push it to ecr
```bash
$ cd provisioning/ecs
$ ./deploy-container-image-to-ecr.sh
```

3. Check
```bash
$ aws ecr list-images --repository-name test | jq '.imageIds | .[].imageTag'
"ecs-serchdomain-sidecar"
"logstash"
```

## Build Lambda
1. Create s3 bucket
```bash
$ s3-ops create lambda-xxxxxxxxxxxxxxxxx-artifact
$ export S3_BUCKET_LAMBDA_NAME="lambda-xxxxxxxxxxxxxxxxx-artifact"
```

2. Build
```bash
$ cd provisioning/lambda
$ ./build.sh
```

3. Push artifact to s3 
```bash
$ s3-ops push $S3_BUCKET_LAMBDA_NAME populate-alb-tg-with-opensearch/populate-alb-tg-with-opensearch.zip
```

## TOOLS
* S3
```bash
$ s3-ops 

Usage: /path/to/ElasticStack/provisioning/helper-scripts/s3-ops.sh MODE ARGS

Mode:     Args:
list      
create    s3-name
delete    s3-name
push      s3-name     zip-file-path
```

* Cloudformation
```bash
$ cfn-stack-ops

Usage: /path/to/ElasticStack/provisioning/helper-scripts/cfn-stack-ops.sh MODE ARGS

Mode:     Args:
create    stack-name s3-bucket [param1=val1 param2=val2]
update    stack-name s3-bucket [param1=val1 param2=val2]
package   path-to-cfn-template-file s3-bucket output-template-file
deploy    stack-name path-to-cfn-template-filee
list      
describe  stack-name
validate  s3-bucket
delete    stack-name
```

# SSH tunnel

```bash
$ sh -i ~/.ssh/your-key.pem ec2-user@your-ec2-instance-public-ip -N -L 9200:vpc-domain-name.region.es.amazonaws.com:443
```