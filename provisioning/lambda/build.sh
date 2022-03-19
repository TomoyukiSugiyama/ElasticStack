#!/bin/sh

GOOS=linux go build populate-alb-tg-with-opensearch.go
zip populate-alb-tg-with-opensearch.zip populate-alb-tg-with-opensearch