package main

import (
	"regexp"
	"sort"
	"strings"
)

// Single-pass, word-boundary replacement: "db" never corrupts "database",
// and an inserted value ("auth" → "authentication") is never re-matched.
func replaceLexicon(text string, lexicon map[string]string) string {
	if len(lexicon) == 0 {
		return text
	}

	keys := make([]string, 0, len(lexicon))
	for key := range lexicon {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return len(keys[i]) > len(keys[j]) })

	quoted := make([]string, len(keys))
	for i, key := range keys {
		quoted[i] = regexp.QuoteMeta(key)
	}
	re, err := regexp.Compile(`(?i)\b(` + strings.Join(quoted, "|") + `)\b`)
	if err != nil {
		return text
	}

	return re.ReplaceAllStringFunc(text, func(match string) string {
		return lexicon[strings.ToLower(match)]
	})
}
