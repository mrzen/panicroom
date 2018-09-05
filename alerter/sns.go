package alerter

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"time"
)

// SNSAlerter - Handle alerts by sending a notification to AWS SNS
type SNSAlerter struct {
	TopicArn string `yaml:"arn"`
	Region string `yaml:"region"`
	client *sns.SNS
}

func (a *SNSAlerter) Alert(e Event) error {

	req := &sns.PublishInput{
		TopicArn: &a.TopicArn,
		Subject: aws.String("[PanicRoom] " + e.WatcherName),
	}

	message := fmt.Sprintf("Path: %s\nOperation:%s\nTimestamp:%s", e.Path, e.Operation, e.T.Format(time.RFC850))
	req.Message = &message

	req.MessageAttributes = map[string]*sns.MessageAttributeValue{
		"operation": {
			StringValue: aws.String(e.Operation.String()),
			DataType: aws.String("String"),
		},
	}

	_, err := a.client.Publish(req)
	return err
}

func (a *SNSAlerter) SetConfig(config interface{}) error {

	m, ok := config.(map[string]interface{})

	if !ok {
		return errors.New("invalid configuration")
	}

	topicArn, ok := m["arn"].(string)

	if !ok || topicArn == "" {
		return errors.New("missing required parameter: arn")
	}

	region, ok := m["region"].(string)

	if !ok || region == "" {
		return errors.New("missing required prameter: region")
	}

	a.TopicArn = topicArn

	s, err := session.NewSession(&aws.Config{
		Region: &region,
	})

	a.client = sns.New(s)

	return err
}