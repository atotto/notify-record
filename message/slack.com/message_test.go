package slack_com

import (
	"testing"

	"github.com/atotto/notify-record/message/internal/notification"
)

func TestMessage(t *testing.T) {
	tests := []struct {
		msg    *notification.Message
		expect string
	}{
		{
			&notification.Message{
				Domain: "example.slack.com",
				Title:  "New notification from Slackbot",
				Body:   "テスト",
			},
			`example.slack.com [Slackbot] テスト`,
		},
		{
			&notification.Message{
				Domain: "example.slack.com",
				Title:  "New notification in #general",
				Body:   "ato: あいうえお",
			},
			`example.slack.com [#general] ato: あいうえお`,
		},
	}

	for n, tt := range tests {
		line := ToString(tt.msg)
		if line != tt.expect {
			t.Fatalf("#%d got %s; want %s", n, line, tt.expect)
		}
	}
}
