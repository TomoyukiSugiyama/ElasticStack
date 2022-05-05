#!/bin/sh

for file in `\find . -name '*.guard'`; do
    filename=`basename $file .guard`
    cat ../provisioning/cfn/$filename.yaml | cfn-guard validate -r $file
done