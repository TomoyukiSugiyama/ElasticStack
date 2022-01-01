#!/bin/sh
# referenced from https://www.slideshare.net/yktko/cloudformation-getting-started-with-yaml

mode=$1; shift
stack_name=$1; shift
template=$1; shift
if [ "$mode" != "create" -a "$mode" != "update" ]; then
    echo "$0 (create|update) stack-name template-name [param1=val1 param2=val2]"; exit 1
fi

params=$(echo $* | perl -pe "s/([^= ]+)=([^ ]+)/ParameterKey=\1,ParameterValue=\2/g")

cmd="aws cloudformation ${mode}-stack --stack-name ${stack_name} --template-body file://${template} --capabilities CAPABILITY_IAM $params"
echo ${cmd}
eval ${cmd}