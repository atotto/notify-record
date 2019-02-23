package calendar_google_com

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
			notification.FromText("lunch", `calendar.google.com

14:00 – 15:00
office`),
			`calendar.google.com [14:00-15:00] > lunch at office`,
		},
		{
			&notification.Message{
				Domain: "calendar.google.com",
				Title:  "lunch",
				Body:   "14:00 – 15:00",
			},
			`calendar.google.com [14:00-15:00] > lunch`,
		},
	}

	for n, tt := range tests {
		line := ToString(tt.msg)
		if line != tt.expect {
			t.Fatalf("#%d got %s; want %s", n, line, tt.expect)
		}
	}
}
