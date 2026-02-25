package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/notify-record/config"
	"github.com/atotto/notify-record/message"
	"github.com/fatih/color"
	"github.com/godbus/dbus"
)

type HookExecution struct {
	Hook    config.Hook
	Message interface{}
}

func getConfigPath() string {
	// 1. Check XDG_CONFIG_HOME
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "notify-record", "config.yml")
	}

	// 2. Use XDG default: ~/.config/notify-record/config.yml
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}
	return filepath.Join(homeDir, ".config", "notify-record", "config.yml")
}

func main() {
	configPath := getConfigPath()

	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("Config loading error: %v", err)
	}

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
	ignoreKeywords := strings.Split(os.Getenv("NOTIFY_RECORD_IGNORE_KEYWORDS"), ",")

	domainColor := color.New(color.FgBlue).SprintFunc()
	headerColor := color.New(color.FgBlack, color.Bold).SprintFunc()
	bodyColor := color.New(color.FgBlack).SprintFunc()
	keywordColor := color.New(color.FgRed, color.Bold).SprintFunc()
	titleColor := color.New(color.FgGreen).SprintFunc()
	tsColor := color.New(color.FgHiRed).SprintFunc()

	hookChan := make(chan HookExecution, 16)

	go func() {
		for hookExec := range hookChan {
			if err := hookExec.Hook.Exec(hookExec.Message); err != nil {
				log.Printf("Hook script execution error: %v", err)
			}
		}
	}()

	var messageString string
	for v := range mch {
		if len(v.Body) == 0 {
			continue
		}
		m := message.Message(v)
		if messageString == m.String() {
			// suppress dup message
			continue
		}
		messageString = m.String()

		// ignore keywords
		if !hasKeywords(messageString, keywords) && hasKeywords(messageString, ignoreKeywords) {
			continue
		}

		if !config.PassesGlobalFilters(m.Domain, m.Title, m.Body) {
			continue
		}

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

		hooks := config.FindHooksForMessage(m.Domain, m.Title, m.Body)
		for _, hook := range hooks {
			hookChan <- HookExecution{Hook: hook, Message: m}
		}
	}
	close(hookChan)
}

func hasKeywords(s string, keys []string) bool {
	for _, key := range keys {
		if strings.Contains(s, key) {
			return true
		}
	}
	return false
}
