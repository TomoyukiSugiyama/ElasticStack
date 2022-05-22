#!/bin/sh

for SERVICE in dummy-log-generator;
do
  (cd ${SERVICE} ; go build ${SERVICE}.go)
done
