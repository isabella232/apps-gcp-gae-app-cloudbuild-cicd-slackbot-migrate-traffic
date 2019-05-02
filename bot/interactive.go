package bot

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/nlopes/slack/slackevents"
	"google.golang.org/api/appengine/v1"
	"google.golang.org/api/option"

	"github.com/qushot/appengine-cloudbuild-slackbot/bot/shared"
)

// 対話用
func Interactive(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body, err := url.QueryUnescape(strings.TrimPrefix(buf.String(), "payload="))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	log.Printf("body: %s", body)

	actionEvent, err := slackevents.ParseActionEvent(body, slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: shared.VerificationToken}))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if actionEvent.Type == "interactive_message" {
		if _, _, err := shared.API.DeleteMessage(shared.ChannelID, actionEvent.MessageTimestamp.String()); err != nil {
			log.Printf("[ERROR] delete message error: %v", err)
			return
		}

		var value string
		if len(actionEvent.Actions) > 0 && len(actionEvent.Actions[0].SelectedOptions) > 0 {
			value = actionEvent.Actions[0].SelectedOptions[0].Value
		}

		if value != "none" {
			if err := shared.PostTextMessage(shared.API, fmt.Sprintf("トラフィックを `%s` に切り替えます！", value)); err != nil {
				return
			}

			go func() {
				scopes := []string{
					appengine.CloudPlatformScope,
				}

				ctx := context.Background()
				svc, err := appengine.NewService(ctx, option.WithScopes(scopes...))
				if err != nil {
					log.Printf("appengine NewService error: %v", err)
					return
				}

				service := &appengine.Service{
					Split: &appengine.TrafficSplit{
						Allocations: map[string]float64{value: 1},
						ShardBy:     "IP",
					},
				}

				if _, err := svc.Apps.Services.Patch(shared.ProjectID, "default", service).MigrateTraffic(false).UpdateMask("split").Do(); err != nil {
					log.Printf("[ERROR] appengine service patch error: %v", err)
					return
				}

				if err := shared.PostTextMessage(shared.API, "切り替えに成功しました！"); err != nil {
					return
				}
			}()

			return
		}

		if err := shared.PostTextMessage(shared.API, "切り替えません"); err != nil {
			return
		}
	}
}
