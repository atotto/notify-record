package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/atotto/notify-record/message"
	"github.com/fatih/color"

	"github.com/godbus/dbus"
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}

	if err = conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, "type='method_call',path='/org/freedesktop/Notifications',member='Notify',eavesdrop='true'").Err; err != nil {
		log.Fatal(err)
	}
	mch := make(chan *dbus.Message, 10)
	conn.Eavesdrop(mch)

	keywords := strings.Split(os.Getenv("NOTIFY_RECORD_KEYWORDS"), ",")

	domainColor := color.New(color.FgBlue).SprintFunc()
	headerColor := color.New(color.FgBlack, color.Bold).SprintFunc()
	bodyColor := color.New(color.FgBlack).SprintFunc()
	keywordColor := color.New(color.FgRed, color.Bold).SprintFunc()
	titleColor := color.New(color.FgGreen).SprintFunc()
	tsColor := color.New(color.FgHiRed).SprintFunc()

	for v := range mch {
		m := message.Message(v)
		body := bodyColor(strings.ReplaceAll(m.Body, "\n", "\n  "))
		if m.Header != "" {
			body = fmt.Sprintf("%s: %s", headerColor(m.Header), body)
		}
		for _, keyword := range keywords {
			if i := strings.Index(body, keyword); i > 0 {
				body = body[:i] + keywordColor(keyword) + body[i+len(keyword):]
			}
		}
		fmt.Printf("%s %s %s\n  %s\n",
			tsColor(m.TimeStamp.Format("15:04:05")),
			domainColor(m.Domain),
			titleColor(m.Title),
			body,
		)
	}
}
