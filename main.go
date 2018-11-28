package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getInstances() (interface{}, error) {
	var region, profile, project string

	flag.StringVar(&region, "region", "ca-central-1", "Region to look for instances - default to ca-central-1")
	flag.StringVar(&profile, "profile", "default", "AWS Config profile to use for call")
	flag.StringVar(&project, "project", "", "Project to look for (via tags)")
	flag.Parse()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:                 profile,
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	}))

	svc := ec2.New(sess, &aws.Config{
		Region: &region,
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
			if len(project) == 0 {
				fmt.Println(*instance.InstanceId)
				fmt.Println(*instance.PublicDnsName)
				fmt.Println("")
			} else {
				for _, tag := range instance.Tags {
					if *tag.Key == "project" && *tag.Value == project {
						fmt.Println(*tag.Value)
						fmt.Println(*instance.InstanceId)
						fmt.Println(*instance.PublicDnsName)
						fmt.Println("")
					}				
				}			
			}
		}
	}

	return res, err
}

func main() {

	if len(os.Args) < 3 {
		panic("You must supply region and profile as arguments!")
	}

	getInstances()
}
