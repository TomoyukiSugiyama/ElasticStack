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

func ResolveIpAddress(domainEndpoint string) (resolvedIpAddress string) {
	ipAddr, err := net.ResolveIPAddr("ip", domainEndpoint)
	if err != nil {
		fmt.Println("Resolve error ", err)
		os.Exit(1)
	}
	resolvedIpAddress = ipAddr.IP.String()
	return
}

func Init() (svc *elasticloadbalancingv2.Client) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		os.Exit(1)
	}

	svc = elasticloadbalancingv2.NewFromConfig(cfg)
	return
}

func GetSpecifiedLoadbalancer(svc *elasticloadbalancingv2.Client, loadBalancerName string) (target types.LoadBalancer) {
	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{}
	resp, err := svc.DescribeLoadBalancers(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to get loadbalancers, %v", err)
		os.Exit(1)
	}

	for _, lb := range resp.LoadBalancers {
		if *lb.LoadBalancerName == loadBalancerName {
			target = lb
			return
		}
	}
	log.Fatalf("failed to get specified loadbalancer")
	os.Exit(1)
	return
}

func GetSpecifiedTargetGroup(svc *elasticloadbalancingv2.Client, loadBalancer types.LoadBalancer, targetGroupName string) (target types.TargetGroup) {
	input := &elasticloadbalancingv2.DescribeTargetGroupsInput{LoadBalancerArn: loadBalancer.LoadBalancerArn}
	resp, err := svc.DescribeTargetGroups(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to get target groups, %v", err)
		os.Exit(1)
	}

	for _, tg := range resp.TargetGroups {
		if *tg.TargetGroupName == targetGroupName {
			target = tg
			return
		}
	}
	log.Fatalf("failed to get specified target group")
	os.Exit(1)
	return
}

func HasTarget(svc *elasticloadbalancingv2.Client, tg types.TargetGroup, ipAddr string) (hasTarget bool) {
	input := &elasticloadbalancingv2.DescribeTargetHealthInput{TargetGroupArn: tg.TargetGroupArn}
	resp, err := svc.DescribeTargetHealth(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to get target helth, %v", err)
		os.Exit(1)
	}

	for _, tgh := range resp.TargetHealthDescriptions {
		if *tgh.Target.Id == ipAddr {
			hasTarget = true
			return
		}
	}
	hasTarget = false
	return
}

func RegisterSpecifiedTarget(svc *elasticloadbalancingv2.Client, tg types.TargetGroup, addr string, port int32) {
	registerTarget := types.TargetDescription{AvailabilityZone: nil, Id: &addr, Port: &port}
	registerTargets := []types.TargetDescription{registerTarget}
	registerTargetInput := &elasticloadbalancingv2.RegisterTargetsInput{TargetGroupArn: tg.TargetGroupArn, Targets: registerTargets}
	_, err := svc.RegisterTargets(context.TODO(), registerTargetInput)
	if err != nil {
		log.Fatalf("failed to register, %v", err)
		os.Exit(1)
	}
}

func DeregisterUnheltyTargets(svc *elasticloadbalancingv2.Client, tg types.TargetGroup) {
	input := &elasticloadbalancingv2.DescribeTargetHealthInput{TargetGroupArn: tg.TargetGroupArn}
	resp, err := svc.DescribeTargetHealth(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to get target helth, %v", err)
		os.Exit(1)
	}

	for _, tgh := range resp.TargetHealthDescriptions {
		if tgh.TargetHealth.State == types.TargetHealthStateEnumUnhealthy {
			DeregisterSpecifiedTarget(svc, tg, *tgh.Target.Id, *tgh.Target.Port)
		}
	}
}

func DeregisterSpecifiedTarget(svc *elasticloadbalancingv2.Client, tg types.TargetGroup, addr string, port int32) {
	deregisterTarget := types.TargetDescription{AvailabilityZone: nil, Id: &addr, Port: &port}
	deregisterTargets := []types.TargetDescription{deregisterTarget}
	deregisterTargetInput := &elasticloadbalancingv2.DeregisterTargetsInput{TargetGroupArn: tg.TargetGroupArn, Targets: deregisterTargets}
	_, err := svc.DeregisterTargets(context.TODO(), deregisterTargetInput)
	if err != nil {
		log.Fatalf("failed to deregister, %v", err)
		os.Exit(1)
	}
}

func HandleLambdaEvent() {
	opensearchAddr := ResolveIpAddress("vpc-my-es-sk5xpobbjxtur7njpsc7qplwlq.ap-northeast-1.es.amazonaws.com")
	fmt.Printf("Opensearch address: %s\n", opensearchAddr)

	svc := Init()

	lb := GetSpecifiedLoadbalancer(svc, "f-iot-alb")
	fmt.Printf("LoadBalancer name : %s\n", *lb.LoadBalancerName)
	fmt.Printf("LoardBalancer arn : %s\n", *lb.LoadBalancerArn)

	tg := GetSpecifiedTargetGroup(svc, lb, "f-iot-alb-tg")
	fmt.Printf("TargetGroup name : %s\n", *tg.TargetGroupName)
	fmt.Printf("TargetGroup arn : %s\n", *tg.TargetGroupArn)

	if !HasTarget(svc, tg, opensearchAddr) {
		const httpsPort = 443
		RegisterSpecifiedTarget(svc, tg, opensearchAddr, httpsPort)
		fmt.Println("Register new target")
	}

	DeregisterUnheltyTargets(svc, tg)
	fmt.Println("Deregister unhealty targets")
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
