package calendar_google_com

import (
	"fmt"
	"strings"

	"github.com/atotto/notify-record/message/internal/notification"
)

func ToString(m *notification.Message) string {
	ss := strings.SplitN(m.Body, "\n", 2)
	if len(ss) < 1 {
		return m.String()
	}

	schedule := strings.Replace(ss[0], " – ", "-", 1)

	if len(ss) == 1 {
		return fmt.Sprintf("%s [%s] > %s", m.Domain, schedule, m.Title)
	}

	location := ss[1]

	return fmt.Sprintf("%s [%s] > %s at %s", m.Domain, schedule, m.Title, location)
}
