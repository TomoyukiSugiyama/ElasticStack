#!/bin/sh

for {SERVICE} in populate-alb-tg-with-opensearch detach-task-to-be-terminated-from-nlb;
do
  (cd ${SERVICE} ; GOOS=linux go build ${SERVICE}.go)
  (cd ${SERVICE} ; zip ${SERVICE}.zip ${SERVICE})
done
