#!/bin/sh

mode=$1; shift
stack_name=$1; shift

if [ "$mode" != "list" -a "$mode" != "describe"]; then
    echo "$0 (list|describe) stack-name"; exit 1
fi

if [ "$mode" == "list" ]; then
params="--stack-status-filter CREATE_COMPLETE"
fi

if [ "$mode" == "describe" ]; then
params="--stack-name ${stack_name}"
fi

cmd="aws cloudformation ${mode}-stacks ${params}"
echo ${cmd}
eval ${cmd}
