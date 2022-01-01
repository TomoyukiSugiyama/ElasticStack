#!/bin/sh
# ref from https://www.slideshare.net/yktko/cloudformation-getting-started-with-yaml

mode=$1; shift
stack_name=$1; shift
template=$1; shift
if [ "$mode" != "create" -a "$mode" != "update" -a "$mode" != "delete" -a "$mode" != "list" -a "$mode" != "describe" ]; then
    echo ""
    echo "Usage: $0 MODE [OPTIONS]"
    echo ""
    echo "Mode      [OPTIONS]"
    echo "create    stack-name cfn-template-name [param1=val1 param2=val2]"
    echo "update    stack-name cfn-template-name [param1=val1 param2=val2]"
    echo "list      stack-name";
    echo "describe  stack-name";
    echo "delete    stack-name"; exit 1
fi

if [ "$mode" == "create" -o "$mode" == "update" ]; then
    params=$(echo $* | perl -pe "s/([^= ]+)=([^ ]+)/ParameterKey=\1,ParameterValue=\2/g")
    options="--template-body file://${template} --capabilities CAPABILITY_IAM ${params}"
    stack_name_option="--stack-name ${stack_name}"
    mode_option="${mode}-stack"
fi

if [ "$mode" = "delete" ]; then
    options=""
    stack_name_option="--stack-name ${stack_name}"
    mode_option="${mode}-stack"
fi

if [ "$mode" == "list" ]; then
    options="--stack-status-filter CREATE_COMPLETE"
    stack_name_option=""
    mode_option="${mode}-stacks"
fi

if [ "$mode" == "describe" ]; then
    options=""
    stack_name_option="--stack-name ${stack_name}"
    mode_option="${mode}-stacks"
fi

cmd="aws cloudformation ${mode_option} ${stack_name_option} ${options}"
echo ${cmd}
eval ${cmd}