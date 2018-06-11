package apitool

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/sirupsen/logrus"
	"time"
)

type SlackAdditional struct {
	Name  string
	Value string
}

type SlackMessage string

func SendLogToSlack(SlackAPILogWebHook string, message SlackMessage, context Context, args ...SlackAdditional) {

	// Send message only if in prod
	if context.AppEnvironment != Prod {
		return
	}

	// Add mandatory attachments
	attachment := slack.Attachment{}
	attachment.AddField(slack.Field{
		Title: "Time",
		Value: time.Now().Format(time.RFC1123),
	}).AddField(slack.Field{
		Title: "Micro Service Name",
		Value: context.ServiceName,
	})

	// Add optional attachments
	for k := range args {
		attachment.AddField(slack.Field{
			Title: args[k].Name,
			Value: args[k].Value,
		})
	}

	// Send payload
	err := slack.Send(SlackAPILogWebHook, "", slack.Payload{
		Text:        string(message),
		Attachments: []slack.Attachment{attachment},
	})

	// If any errors, log them in DB
	if len(err) > 0 {
		logrus.Errorf("Could not send log to slack\n%v", err)
	}
}
