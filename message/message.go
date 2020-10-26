package message

import (
	"strings"

	"github.com/atotto/notify-record/message/internal/notification"

	calendar_google_com "github.com/atotto/notify-record/message/calendar.google.com"
	slack_com "github.com/atotto/notify-record/message/slack.com"

	"github.com/godbus/dbus"
)

func String(m *dbus.Message) string {
	msg := notification.From(m)
	switch {
	case strings.HasSuffix(msg.Domain, ".slack.com"):
		return slack_com.ToString(msg)
	case msg.Domain == "calendar.google.com":
		return calendar_google_com.ToString(msg)
	default:
		return msg.String()
	}
}

func Message(m *dbus.Message) *notification.Message {
	msg := notification.From(m)
	switch {
	case strings.HasSuffix(msg.Domain, ".slack.com"):
		return slack_com.Message(msg)
	case msg.Domain == "calendar.google.com":
		return calendar_google_com.Message(msg)
	default:
		return msg
	}
}
