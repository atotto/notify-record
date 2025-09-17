package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Hooks   []Hook   `yaml:"hooks"`
	Filters []Filter `yaml:"filter"`
}

type Filter struct {
	Domain string               `yaml:"domain"`
	Title  TitleFilter          `yaml:"title"`
	Body   BodyFilterWithLength `yaml:"body"`
}

type TitleFilter struct {
	Excludes []string `yaml:"excludes"`
}

type BodyFilterWithLength struct {
	Contains []string    `yaml:"contains"`
	Excludes []string    `yaml:"excludes"`
	Length   LengthLimit `yaml:"length"`
}

type LengthLimit struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		return &Config{}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

func (config *Config) FindHooksForMessage(domain, title, body string) []Hook {
	if config == nil {
		return nil
	}
	var hooks []Hook
	for _, hook := range config.Hooks {
		if hook.MatchesMessage(domain, title, body) {
			hooks = append(hooks, hook)
		}
	}
	return hooks
}

func (config *Config) PassesGlobalFilters(domain, title, body string) bool {
	if config == nil || len(config.Filters) == 0 {
		return true
	}

	for _, filter := range config.Filters {
		if filter.MatchesMessage(domain, title, body) {
			return false
		}
	}
	return true
}

func matchesDomain(pattern, domain string) bool {
	if pattern == domain {
		return true
	}

	matched, err := filepath.Match(pattern, domain)
	if err != nil {
		return false
	}
	return matched
}
