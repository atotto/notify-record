package notification

import "testing"

func TestMessageParseBody(t *testing.T) {
	body := `example.slack.com

テスト`

	domain, msg := parseBody(body)

	if domain != "example.slack.com" {
		t.Fatalf("got %s; want example.com", domain)
	}
	if msg != "テスト" {
		t.Fatalf("got %s; want テスト", msg)
	}
}
