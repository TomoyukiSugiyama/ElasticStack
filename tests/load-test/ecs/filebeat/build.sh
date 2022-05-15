#!/bin/sh

for SERVICE in dummy-log-generator;
do
  (cd ${SERVICE} ; GOOS=linux go build ${SERVICE}.go)
done
