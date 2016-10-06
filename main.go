package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/certifi/gocertifi"
	"github.com/kr/pretty"
)

func main() {
	err := do()
	if err != nil {
		log.Fatal(err)
	}
}

func do() error {
	cert_pool, err := gocertifi.CACerts()
	if err != nil {
		return err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: cert_pool},
		},
	}

	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	metadata := ec2metadata.New(sess, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody))

	info, err := metadata.GetInstanceIdentityDocument()
	if err != nil {
		return err
	}

	conf := aws.NewConfig().
		WithRegion(info.Region).
		WithHTTPClient(client).
		WithLogLevel(aws.LogDebugWithHTTPBody)

	ec2svc := ec2.New(sess, conf)

	resp, err := ec2svc.ModifyInstanceAttribute(&ec2.ModifyInstanceAttributeInput{
		Attribute: aws.String(ec2.InstanceAttributeNameSourceDestCheck),
		SourceDestCheck: &ec2.AttributeBooleanValue{
			Value: aws.Bool(false),
		},
	})
	if err != nil {
		return err
	}

	pretty.Log(resp)
	return nil
}
