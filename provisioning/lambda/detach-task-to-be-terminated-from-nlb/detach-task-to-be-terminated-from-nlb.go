// main TerminationNotificationを受け取ったECSタスクを、NLBのターゲットから削除します。
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

func Init() (svc *elasticloadbalancingv2.Client) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK config, %v", err)
		os.Exit(1)
	}

	svc = elasticloadbalancingv2.NewFromConfig(cfg)
	return
}

func GetSpecifiedLoadbalancer(svc *elasticloadbalancingv2.Client, loadBalancerId string) (target types.LoadBalancer) {
	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{}
	resp, err := svc.DescribeLoadBalancers(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to get loadbalancers, %v", err)
		os.Exit(1)
	}

	for _, lb := range resp.LoadBalancers {
		if *lb.LoadBalancerArn == loadBalancerId {
			target = lb
			return
		}
	}
	log.Fatalf("failed to get specified loadbalancer")
	os.Exit(1)
	return
}

func GetSpecifiedTargetGroup(svc *elasticloadbalancingv2.Client, loadBalancer types.LoadBalancer, targetGroupId string) (target types.TargetGroup) {
	input := &elasticloadbalancingv2.DescribeTargetGroupsInput{LoadBalancerArn: loadBalancer.LoadBalancerArn}
	resp, err := svc.DescribeTargetGroups(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to get target groups, %v", err)
		os.Exit(1)
	}

	for _, tg := range resp.TargetGroups {
		if *tg.TargetGroupArn == targetGroupId {
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

type NetworkInterface struct {
	PrivateIpv4Address string `json:"privateIpv4Address"`
}
type Container struct {
	NetworkInterfaces []NetworkInterface
	Name              string `json:"name"`
}
type EcsEvent struct {
	StopCode   string `json:"stopCode"`
	Containers []Container
}

func HandleLambdaEvent(_ context.Context, event events.CloudWatchEvent) {
	var ecsEvent EcsEvent
	if err := json.Unmarshal(event.Detail, &ecsEvent); err != nil {
		os.Exit(1)
	}

	fmt.Printf("stopCode = %s\n", ecsEvent.StopCode)
	if ecsEvent.StopCode != "TerminationNotice" {
		return
	}
	var ecsIp string
	for _, contaier := range ecsEvent.Containers {
		if contaier.Name == "logstash" {
			for _, ni := range contaier.NetworkInterfaces {
				ecsIp = ni.PrivateIpv4Address
			}
		}
	}
	fmt.Printf("ip v4 = %s\n", ecsIp)

	nlbId := os.Getenv("NlbId")
	nlbTargetGroupId := os.Getenv("NlbTargetGroupId")
	fmt.Printf("GET ENV AlbId: %s AlbTargetGroupId: %s\n", nlbId, nlbTargetGroupId)

	svc := Init()

	lb := GetSpecifiedLoadbalancer(svc, nlbId)
	fmt.Printf("GET LoadbalancerName: %s LoadbalancerArn: %s\n", *lb.LoadBalancerName, *lb.LoadBalancerArn)

	tg := GetSpecifiedTargetGroup(svc, lb, nlbTargetGroupId)
	fmt.Printf("GET TargetGroupName: %s TargetGroupArn: %s\n", *tg.TargetGroupName, *tg.TargetGroupArn)

	if HasTarget(svc, tg, ecsIp) {
		const tcpPort = 5044
		DeregisterSpecifiedTarget(svc, tg, ecsIp, tcpPort)
		fmt.Println("DEREGISTER")
	}
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
