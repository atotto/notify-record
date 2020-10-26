package slack_com

import (
	"fmt"
	"strings"

	"github.com/atotto/notify-record/message/internal/notification"
)

func Message(m *notification.Message) *notification.Message {
	index := strings.LastIndex(m.Title, " ")
	title := m.Title[index+1:]
	m.Title = title
	return m
}

func ToString(m *notification.Message) string {
	index := strings.LastIndex(m.Title, " ")
	title := m.Title[index+1:]
	return fmt.Sprintf("%s [%s] %s", m.Domain, title, m.Body)
}
