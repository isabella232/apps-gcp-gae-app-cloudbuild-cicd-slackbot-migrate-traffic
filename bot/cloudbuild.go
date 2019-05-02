package bot

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nlopes/slack"
	"google.golang.org/api/appengine/v1"
	"google.golang.org/api/option"

	"github.com/qushot/appengine-cloudbuild-slackbot/bot/shared"
)

// PubSubのメッセージを受け取る
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// PubSubMessageパース用
type (
	build struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		LogURL string `json:"logUrl"`
		Source source `json:"source"`
	}

	repoSource struct {
		RepoName   string `json:"repoName"`
		BranchName string `json:"branchName"`
	}

	source struct {
		RepoSource repoSource `json:"repoSource"`
	}
)

// Cloud Buildの結果通知用
func CloudBuild(ctx context.Context, m PubSubMessage) error {
	b := new(build)
	if err := json.Unmarshal(m.Data, b); err != nil {
		log.Printf("[ERROR] %v", err)
		return err
	}

	if b.Status == "SUCCESS" && b.Source.RepoSource.BranchName != "" {
		log.Printf("PubSubMessage: %s", m.Data)

		if err := shared.PostTextMessage(shared.API, "ビルド成功"); err != nil {
			return err
		}

		attachments, err := newSuccessMsgAttachments(shared.ProjectID)
		if err != nil {
			return err
		}

		if err := shared.PostAttachmentMessage(shared.API, newSuccessAttachmentOptions(attachments)); err != nil {
			return err
		}

	} else if (b.Status == "FAILURE" || b.Status == "INTERNAL_ERROR" || b.Status == "TIMEOUT") && b.Source.RepoSource.BranchName != "" {
		log.Printf("PubSubMessage: %s", m.Data)

		if err := shared.PostTextMessage(shared.API, "ビルド失敗"); err != nil {
			return err
		}
	}

	return nil
}

func newSuccessAttachmentOptions(attachActOpts []slack.AttachmentActionOption) *[]slack.MsgOption {
	return &[]slack.MsgOption{
		slack.MsgOptionAttachments(slack.Attachment{
			Color:      "good",
			Fallback:   "select version",
			CallbackID: "version_selection",
			Text:       "どのバージョンに切り替える？",
			Actions: []slack.AttachmentAction{
				{
					Name:    "version_list",
					Text:    "Version List",
					Type:    "select",
					Options: attachActOpts,
				},
			},
			MarkdownIn: []string{"text"},
		}),
	}
}

func newSuccessMsgAttachments(pjID string) ([]slack.AttachmentActionOption, error) {
	ctx := context.Background()
	svc, err := appengine.NewService(ctx, option.WithScopes([]string{
		appengine.AppengineAdminScope,
	}...))
	if err != nil {
		log.Printf("[ERROR] appengine NewService error: %v", err)
		return nil, err
	}

	cv, err := svc.Apps.Services.Get(pjID, "default").Do()
	if err != nil {
		log.Printf("[ERROR] appengine service patch error: %v", err)
		return nil, err
	}

	var version string
	for k, _ := range cv.Split.Allocations {
		version = k
		log.Printf("Current version: %s", k)
	}

	list, err := svc.Apps.Services.Versions.List(pjID, "default").View("FULL").Do()
	if err != nil {
		log.Printf("[ERROR] appengine service version list error: %v", err)
		return nil, err
	}

	opts := make([]slack.AttachmentActionOption, len(list.Versions)+1)
	opts[0].Text = "切り替えない"
	opts[0].Value = "none"

	for i, v := range list.Versions {
		opts[i+1].Text = v.Id
		opts[i+1].Value = v.Id
		if version == v.Id {
			opts[i+1].Description = "default"
		}
	}

	return opts, nil
}
