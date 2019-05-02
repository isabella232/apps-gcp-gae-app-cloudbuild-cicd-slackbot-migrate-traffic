package shared

import (
	"log"
	"os"

	"github.com/nlopes/slack"
)

var (
	API               = slack.New(os.Getenv("BOT_USER_OAUTH_ACCESS_TOKEN"))
	ChannelID         = os.Getenv("CHANNEL_ID")
	ProjectID         = os.Getenv("GCP_PROJECT")
	VerificationToken = os.Getenv("VERIFICATION_TOKEN")
)

func PostTextMessage(api *slack.Client, text string) error {
	opt := &[]slack.MsgOption{
		slack.MsgOptionText(text, false),
	}

	return postMessage(api, opt)
}

func PostAttachmentMessage(api *slack.Client, opt *[]slack.MsgOption) error {
	return postMessage(api, opt)
}

func postMessage(api *slack.Client, opts *[]slack.MsgOption) error {
	_, ts, err := api.PostMessage(ChannelID, *opts...)
	if err != nil {
		log.Printf("[ERROR] post message error: %v", err)
		return err
	}

	log.Printf("success post messsage. (TimeStamp: %s)", ts)
	return nil
}
