package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var roleArn = os.Getenv("AWS_PROFILE_ARN")

func getRegion() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter region: ")

	input, _ := reader.ReadString('\n')

	region := strings.TrimRight(input, "\n")

	return region
}

func getInstances() (interface{}, error) {
	region := getRegion()

	sess := session.Must(session.NewSession())

	creds := stscreds.NewCredentials(sess, roleArn)

	svc := ec2.New(sess, &aws.Config{
		Credentials: creds,
		Region:      &region,
	})

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
					aws.String("pending"),
				},
			},
		},
	}

	res, err := svc.DescribeInstances(params)

	if err != nil {
		log.Fatal(err)
	}

	for _, reservation := range res.Reservations {
		for _, instance := range reservation.Instances {
			for _, tag := range instance.Tags {
				if *tag.Key == "project" {
					fmt.Println(*tag.Value)
				}
			}
			fmt.Println(*instance.InstanceId)
			fmt.Println(*instance.PublicDnsName)
			fmt.Println("")

		}
	}

	return res, err
}

func main() {
	getInstances()
}
