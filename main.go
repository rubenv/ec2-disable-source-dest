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

	metadata := ec2metadata.New(sess)

	info, err := metadata.GetInstanceIdentityDocument()
	if err != nil {
		return err
	}

	conf := aws.NewConfig().
		WithRegion(info.Region).
		WithHTTPClient(client)
	ec2svc := ec2.New(sess, conf)

	_, err = ec2svc.ModifyInstanceAttribute(&ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(info.InstanceID),
		SourceDestCheck: &ec2.AttributeBooleanValue{
			Value: aws.Bool(false),
		},
	})
	return err
}
