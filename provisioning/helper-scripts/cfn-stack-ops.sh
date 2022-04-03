#!/bin/sh
# ref: https://www.slideshare.net/yktko/cloudformation-getting-started-with-yaml

mode=$1; shift
arg1=$1; shift
arg2=$1; shift
arg3=$1; shift
if [ "$mode" != "create" -a "$mode" != "update" -a "$mode" != "delete" -a "$mode" != "list" -a "$mode" != "describe" -a "$mode" != "validate" -a "$mode" != "package"  -a "$mode" != "deploy" ]; then
    echo ""
    echo "Usage: $0 MODE ARGS"
    echo ""
    echo "Mode:     Args:"
    echo "create    stack-name s3-bucket [param1=val1 param2=val2]"
    echo "update    stack-name s3-bucket [param1=val1 param2=val2]"
    echo "package   path-to-cfn-template-file s3-bucket output-template-file"
    echo "deploy    stack-name path-to-cfn-template-filee"
    echo "list      ";
    echo "describe  stack-name";
    echo "validate  s3-bucket"
    echo "delete    stack-name"; exit 1
fi

if [ "$mode" == "create" -o "$mode" == "update" ]; then
    params=$(echo $* | perl -pe "s/([^= ]+)=([^ ]+)/ParameterKey=\1,ParameterValue=\2/g")
    args="--template-url https://s3.amazonaws.com/${arg2}/artifact.yaml --capabilities CAPABILITY_IAM --capabilities CAPABILITY_NAMED_IAM ${params}"
    stack_name_option="--stack-name ${arg1}"
    mode_option="${mode}-stack"
fi

if [ "$mode" == "package" ]; then
    args="--template-file ${arg1} --s3-bucket ${arg2} --output-template-file ${arg3}"
    stack_name_option=""
    mode_option="${mode}"
fi

if [ "$mode" == "deploy" ]; then
    args="--template-file ${arg2} --capabilities CAPABILITY_IAM --capabilities CAPABILITY_NAMED_IAM"
    stack_name_option="--stack-name ${arg1}"
    mode_option="${mode}"
fi

if [ "$mode" == "list" ]; then
    args="--stack-status-filter CREATE_COMPLETE"
    stack_name_option=""
    mode_option="${mode}-stacks"
fi

if [ "$mode" == "describe" ]; then
    args=""
    stack_name_option="--stack-name ${arg1}"
    mode_option="${mode}-stacks"
fi

if [ "$mode" = "validate" ]; then
    args="--template-url https://s3.amazonaws.com/${arg1}/artifact.yaml"
    stack_name_option=""
    mode_option="${mode}-template"
fi

if [ "$mode" = "delete" ]; then
    args=""
    stack_name_option="--stack-name ${arg1}"
    mode_option="${mode}-stack"
fi
cmd="aws cloudformation ${mode_option} ${stack_name_option} ${args}"
echo ${cmd}
eval ${cmd}