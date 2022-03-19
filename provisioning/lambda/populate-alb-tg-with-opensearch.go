package main

import (
	"fmt"
	"net"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleLambdaEvent() {
	addr, err := net.ResolveIPAddr("ip", "vpc-my-es-sk5xpobbjxtur7njpsc7qplwlq.ap-northeast-1.es.amazonaws.com")
	if err != nil {
		fmt.Println("Resolve error ", err)
		os.Exit(1)
	}
	fmt.Println("Addr: ", addr)
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
