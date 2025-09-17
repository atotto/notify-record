package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type BodyFilter struct {
	Contains []string `yaml:"contains"`
	Excludes []string `yaml:"excludes"`
}

type Hook struct {
	Domain string     `yaml:"domain"`
	Title  string     `yaml:"title"`
	Body   BodyFilter `yaml:"body"`
	Script string     `yaml:"script"`
}

func (h Hook) Exec(message interface{}) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message to JSON: %v", err)
	}

	cmd := exec.Command("sh", "-c", h.Script)
	cmd.Stdin = strings.NewReader(string(messageJSON))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (h Hook) MatchesMessage(domain, title, body string) bool {
	// Domain matching
	if !matchesDomain(h.Domain, domain) {
		return false
	}

	// Title matching
	if h.Title != "" && !strings.Contains(title, h.Title) {
		return false
	}

	// Body filtering
	if !h.matchesBodyFilter(body) {
		return false
	}

	return true
}

func (h Hook) matchesBodyFilter(body string) bool {
	// Check contains conditions
	if len(h.Body.Contains) > 0 {
		hasRequired := false
		for _, keyword := range h.Body.Contains {
			if strings.Contains(body, keyword) {
				hasRequired = true
				break
			}
		}
		if !hasRequired {
			return false
		}
	}

	// Check excludes conditions
	for _, keyword := range h.Body.Excludes {
		if strings.Contains(body, keyword) {
			return false
		}
	}

	return true
}
