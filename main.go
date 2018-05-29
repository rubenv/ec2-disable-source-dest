package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func die(err error) {
	fmt.Fprintln(os.Stderr, "disable-check:", err)
	os.Exit(1)
}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		die(err)
	}

	metadata := ec2metadata.New(sess)

	info, err := metadata.GetInstanceIdentityDocument()
	if err != nil {
		die(err)
	}

	ec2svc := ec2.New(sess, aws.NewConfig().WithRegion(info.Region))

	_, err = ec2svc.ModifyInstanceAttribute(&ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(info.InstanceID),
		SourceDestCheck: &ec2.AttributeBooleanValue{
			Value: aws.Bool(false),
		},
	})
	if err != nil {
		die(err)
	}
}
