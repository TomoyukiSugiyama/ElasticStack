package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

func HandleLambdaEvent() {
	addr, err := net.ResolveIPAddr("ip", "vpc-my-es-sk5xpobbjxtur7njpsc7qplwlq.ap-northeast-1.es.amazonaws.com")
	if err != nil {
		fmt.Println("Resolve error ", err)
		os.Exit(1)
	}
	fmt.Println("Addr: ", addr)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	fmt.Println(cfg)

	svc := elasticloadbalancingv2.NewFromConfig(cfg)
	//fmt.Println(svc)

	//lbarn := "arn:aws:elasticloadbalancing:ap-northeast-1:645402554699:loadbalancer/app/opensearch-lb/0acadb2211bef59e"

	//input := &elasticloadbalancingv2.DescribeLoadBalancersInput{LoadBalancerArns: []string{lbarn}}
	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{}
	//fmt.Println(input)
	resp, err := svc.DescribeLoadBalancers(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to get loadbalancers, %v", err)
	}

	//fmt.Println(resp)
	for _, lb := range resp.LoadBalancers {
		fmt.Printf("ARN : %s\n", *lb.LoadBalancerArn)
		fmt.Printf("DNS name : %s\n", *lb.DNSName)
		fmt.Printf("LoadBalancer name : %s\n", *lb.LoadBalancerName)
		inputtg := &elasticloadbalancingv2.DescribeTargetGroupsInput{LoadBalancerArn: lb.LoadBalancerArn}
		tgs, err := svc.DescribeTargetGroups(context.TODO(), inputtg)
		if err != nil {
			log.Fatalf("failed to get target groups, %v", err)
		}
		for _, tg := range tgs.TargetGroups {
			fmt.Printf("TargetGroupe name : %s\n", *tg.TargetGroupName)
			fmt.Printf("TargetGroupe ARN : %s\n", *tg.TargetGroupArn)
			fmt.Printf("TargetGroupe port : %d\n", *tg.Port)

		}

	}
}

func main() {
	runtime.Start(HandleLambdaEvent)
}
