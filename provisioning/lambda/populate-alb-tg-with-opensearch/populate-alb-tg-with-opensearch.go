package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
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

	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{}
	resp, err := svc.DescribeLoadBalancers(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to get loadbalancers, %v", err)
	}

	var lbTargetIp string

	var newTarget *types.TargetGroup
	var outputTarget *elasticloadbalancingv2.DescribeTargetHealthOutput
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
			inputhelth := &elasticloadbalancingv2.DescribeTargetHealthInput{TargetGroupArn: tg.TargetGroupArn}

			result, err := svc.DescribeTargetHealth(context.TODO(), inputhelth)
			newTarget = &tg
			outputTarget = result
			if err != nil {
				log.Fatalf("failed to get target helth, %v", err)
			}
			for _, tgh := range result.TargetHealthDescriptions {
				fmt.Printf("Target Id : %s\n", *tgh.Target.Id)
				lbTargetIp = *tgh.Target.Id
			}
		}
	}

	if lbTargetIp == addr.IP.String() {

		fmt.Println("Same Ip")
	} else {

		id := addr.IP.String()
		outputTarget.TargetHealthDescriptions[0].Target.Id = &id
		outputTarget.TargetHealthDescriptions[0].Target.AvailabilityZone = nil
		targets := []types.TargetDescription{*outputTarget.TargetHealthDescriptions[0].Target}
		fmt.Printf("Target Arn : %s\n", *newTarget.TargetGroupArn)
		tginput := &elasticloadbalancingv2.RegisterTargetsInput{TargetGroupArn: newTarget.TargetGroupArn, Targets: targets}
		_, err := svc.RegisterTargets(context.TODO(), tginput)
		if err != nil {
			log.Fatalf("failed to register, %v", err)
		}
		fmt.Println("Defferent Ip")
	}
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
