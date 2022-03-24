package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
)

func echoResource(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
	if event.RequestType == "Delete" {
		return
	}

	addr, err := net.ResolveIPAddr("ip", "vpc-my-es-sk5xpobbjxtur7njpsc7qplwlq.ap-northeast-1.es.amazonaws.com")
	if err != nil {
		fmt.Println("Resolve error ", err)
		os.Exit(1)
	}

	data = map[string]interface{}{
		"Addr": addr.IP.String(),
	}

	response := cfn.NewResponse(&event)
	response.Data = data

	response.Send()

	return
}

func main() {
	lambda.Start(cfn.LambdaWrap(echoResource))
}
