package config

import (
	"strings"
)

func (f Filter) MatchesMessage(domain, title, body string) bool {
	if f.Domain != "" && !matchesDomain(f.Domain, domain) {
		return false
	}
	if f.matchesTitleFilter(title) && f.matchesBodyLengthFilter(body) && f.matchesBodyFilter(body) {
		return true
	}
	return false
}

func (f Filter) matchesTitleFilter(title string) bool {
	for _, keyword := range f.Title.Excludes {
		if strings.Contains(title, keyword) {
			return true
		}
	}
	return false
}

func (f Filter) matchesBodyLengthFilter(body string) bool {
	if f.Body.Length.Min > 0 && len(body) < f.Body.Length.Min {
		return true
	}
	if f.Body.Length.Max > 0 && len(body) > f.Body.Length.Max {
		return true
	}
	return false
}

func (f Filter) matchesBodyFilter(body string) bool {
	for _, keyword := range f.Body.Excludes {
		if strings.Contains(body, keyword) {
			return true
		}
	}
	return false
}
