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
	if event.RequestType == cfn.RequestDelete || event.RequestType == cfn.RequestUpdate {
		return
	}

	domainEndpoint, _ := event.ResourceProperties["DomainEndpoint"].(string)
	addr, err := net.ResolveIPAddr("ip", domainEndpoint)
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
