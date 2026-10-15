package main

import (
	"reflect"
	"testing"
)

func TestSplitSentences(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"Build passed. Tests green.", []string{"Build passed.", "Tests green."}},
		// one-word sentences fold into their neighbor: fewer afplay restarts,
		// and kokoro already pauses inside a chunk
		{"One! Two? Three.", []string{"One! Two? Three."}},
		{"Released v1.2.3 today.", []string{"Released v1.2.3 today."}},
		{"Deadlock in auth_service.go fixed.", []string{"Deadlock in auth_service.go fixed."}},
		{"Hi. Deploy is done now.", []string{"Hi. Deploy is done now."}},
		{"No trailing punctuation", []string{"No trailing punctuation"}},
		{"", nil},
	}
	for _, c := range cases {
		got := SplitSentences(c.in)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("SplitSentences(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestCacheKey(t *testing.T) {
	a := CacheKey("af_heart", 1.0, "hello")
	if a != CacheKey("af_heart", 1.0, "hello") {
		t.Error("same inputs should give the same key")
	}
	if a == CacheKey("af_bella", 1.0, "hello") {
		t.Error("different voice should give a different key")
	}
	if a == CacheKey("af_heart", 1.1, "hello") {
		t.Error("different speed should give a different key")
	}
	if a == CacheKey("af_heart", 1.0, "hello world") {
		t.Error("different text should give a different key")
	}
}
