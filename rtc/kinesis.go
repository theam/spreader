package rtc

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"log"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

const StreamName = "spreader-test"

func NewKinesis() *kinesis.Kinesis{
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
		Credentials: credentials.NewSharedCredentials("", ""),
	})

	if err != nil {
		log.Fatal(err)
	}

	return kinesis.New(sess)
}
