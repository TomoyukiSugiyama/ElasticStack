#!/bin/sh

mode=$1; shift
arg1=$1; shift

if [ "$mode" != "list" -a "$mode" != "create" -a "$mode" != "delete" ]; then
    echo ""
    echo "Usage: $0 MODE ARGS"
    echo ""
    echo "Mode:     Args:"
    echo "list      "
    echo "create    s3-name"
    echo "delete    s3-name"; exit 1
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

cmd="aws s3 ${mode_option} ${args}"
echo ${cmd}
eval ${cmd}