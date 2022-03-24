#!/bin/sh

for SERVICE in populate-alb-tg-with-opensearch populate-alb-tg-with-opensearch-on-cfn;
do
  (cd $SERVICE ; GOOS=linux go build $SERVICE.go)
  (cd $SERVICE ; zip $SERVICE.zip $SERVICE)
done
