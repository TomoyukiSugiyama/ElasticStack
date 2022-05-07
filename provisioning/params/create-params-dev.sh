#!/bin/sh

for file in `\find . -name 'dev-*.yaml'`; do
    aws ssm put-parameter --cli-input-yaml file://${file}
done
