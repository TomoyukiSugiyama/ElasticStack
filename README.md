# ElasticStack

![](https://codebuild.ap-northeast-1.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoickNjQlcrckNSS0xkUGJPUFVkWUdBSlJVcmZwc2NnTWZld2tZVEFLZ2pCR1E3bEdHV1VmS0plYnFKNWJvYmRWeWErSDUrc2hNeERxYTB1RllxOGpvY0E0PSIsIml2UGFyYW1ldGVyU3BlYyI6ImZhU1FGQjNCUkRIbnJST0YiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

## Environment setup
1. Create parameter files to suit your environment.

```
$ cp params/template.yaml params/dev-xxxxx.yaml
```

You need to create following parameters.

| Name                              | Value                                                                                                                                                                                                                                                                        |
| --------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| /dev/github/ConnectionId          | codestar-connections arn to build code from github repogitory. ( Need to connect GitHub with CodeBuild at developper tools console. After that you can get arn from developper tools console. ex. arn:aws:codestar-connections:region:account-id:connection/XXXXXXXXXXXXXX ) |
| /dev/s3/cfn/BucketName            | s3 bucket name for cloudformation files ( This value is used to create s3 bucket for cfn.　Need to check tha rule of s3 bucket name. ex. dev-cloudformation-XXXXXXXXXXXXXXXXXXXX-artifact )                                                                                  |
| /dev/s3/lambda/BucketName         | s3 bucket name for lambda function .zip files ( This value is used to create s3 bucket for cfn.　Need to check tha rule of s3 bucket name. ex. dev-lambda-XXXXXXXXXXXXXXXXXXXX-artifact )                                                                                    |
| /dev/s3/PrefixListId              | managed prefix list id for s3 gateway endpoint ( Get id from vpc console. ex. pl-12a34567 )                                                                                                                                                                                  |
| /dev/slack/codepipeline/ChannelId | slack channel id to notificate result of codepipeline action ( Get id from slack channel. Right click the channel and copy link. ex. A1B2C3D45EF )                                                                                                                           |
| /dev/slack/guardduty/ChannelId    | slack channel id to notificate result from guardduty ( Get id from slack channel. Right click the channel and copy link. ex. A1B2C3D45EF )                                                                                                                                   |
| /dev/slack/WorkspaceId            | slack workspace id to notificate ( Need to connect slack with chat bot at chat bot console. After that you can get id from console. ex. A1B2C3D45EF )                                                                                                                        |
2. Run setup script.

Please install aws cli before run script. Since this script use aws cli.

```
$ ./pre-build/setup.sh
```

## Local Setting (manually)
If you want to deploy manually, please setup local environment.

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
$ cfn-stack-ops package provisioning/cfn/elastic-stack.yaml $S3_BUCKET_CFN_NAME provisioning/artifacts/artifact.yaml
aws cloudformation package --template-file provisioning/cfn/elastic-stack.yaml --s3-bucket elastic-stack-xxxxxxxxxxxxxxxxx-artifact --output-template-file provisioning/artifacts/artifact.yaml
Uploading to XXXXXXXXXXXXXXXXXXXXXXXXXXXX.template  1669 / 1669.0  (100.00%)
Successfully packaged artifacts and wrote output template to file provisioning/artifacts/artifact.yaml.
Execute the following command to deploy the packaged template
aws cloudformation deploy --template-file /path/to/ElasticStack/provisioning/artifacts/artifact.yaml --stack-name <YOUR STACK NAME>
```

3. Deploy cfn

Create stack from local file.
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