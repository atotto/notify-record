package notification

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
)

type Message struct {
	Domain string
	Title  string
	Body   string
}

// https://developer.gnome.org/notification-spec/
type DBusMessage struct {
	Name    string
	ID      uint32
	Icon    string
	Summary string
	Body    string
	Actions []string
	Hints   map[string]dbus.Variant
	Expire  int32
}

func From(m *dbus.Message) *Message {
	msg := DBusMessage{
		Name:    m.Body[0].(string),
		ID:      m.Body[1].(uint32),
		Icon:    m.Body[2].(string),
		Summary: m.Body[3].(string),
		Body:    m.Body[4].(string),
		Actions: m.Body[5].([]string),
		Hints:   m.Body[6].(map[string]dbus.Variant),
		Expire:  m.Body[7].(int32),
	}

	domain, body := parseBody(msg.Body)
	return &Message{
		Domain: domain,
		Title:  msg.Summary,
		Body:   body,
	}
}

func FromText(title, body string) *Message {
	domain, body := parseBody(body)
	return &Message{
		Domain: domain,
		Title:  title,
		Body:   body,
	}
}

func parseBody(body string) (domain string, msg string) {
	ss := strings.SplitN(body, "\n", 3)
	if len(ss) == 3 {
		return ss[0], ss[2]
	}
	return "", body
}

func (m *Message) String() string {
	if m.Domain == "" {
		return fmt.Sprintf("%s: %s", m.Title, m.Body)
	}
	return fmt.Sprintf("%s %s: %s", m.Domain, m.Title, m.Body)
}
