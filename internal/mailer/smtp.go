package mailer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SmtpMailer struct {
	fromEmail string
	client    *ses.Client
}

func aws(s string) *string {
	return &s
}

func NewSESClient(fromEmail, region string) (*SmtpMailer, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}
	return &SmtpMailer{
		client:    ses.NewFromConfig(cfg),
		fromEmail: fromEmail,
	}, nil
}

func (m *SmtpMailer) Send(templateFile, username, email string, data any) error {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}
	input := &ses.SendTemplatedEmailInput{
		Source: &m.fromEmail,
		Destination: &types.Destination{
			ToAddresses: []string{email},
		},
		Template:     aws(templateFile),
		TemplateData: aws(string(dataJson)),
	}
	// retries
	var retryErr error
	for i := range maxRetries {
		_, retryErr = m.client.SendTemplatedEmail(context.TODO(), input)
		if retryErr != nil {
			// exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return nil
	}
	return fmt.Errorf("failed to send email after maxRetries: %v", retryErr.Error())
}
