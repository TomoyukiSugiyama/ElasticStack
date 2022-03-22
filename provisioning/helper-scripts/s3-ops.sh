#!/bin/sh

mode=$1; shift
arg1=$1; shift
arg2=$1; shift

if [ "$mode" != "list" -a "$mode" != "create" -a "$mode" != "delete"  -a "$mode" != "push" ]; then
    echo ""
    echo "Usage: $0 MODE ARGS"
    echo ""
    echo "Mode:     Args:"
    echo "list      "
    echo "create    s3-name"
    echo "delete    s3-name"
    echo "push      s3-name     zip-file-path"; exit 1
fi

if [ "$mode" == "list" ]; then
    args=""
    mode_option="ls"
fi

if [ "$mode" == "create" ]; then
    args="s3://${arg1}"
    mode_option="mb"
fi

if [ "$mode" == "delete" ]; then
    args="s3://${arg1}"
    mode_option="rb"
fi

if [ "$mode" == "push" ]; then
    args="s3://${arg1}"
    mode_option="cp ${arg2}"
fi

cmd="aws s3 ${mode_option} ${args}"
echo ${cmd}
eval ${cmd}