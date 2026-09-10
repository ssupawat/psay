package main

import "testing"

func TestReplaceLexicon(t *testing.T) {
	lex := map[string]string{
		"db":   "database",
		"nil":  "null",
		"auth": "authentication",
	}
	cases := []struct{ in, want string }{
		{"the db is slow", "the database is slow"},
		{"The DB is slow", "The database is slow"},
		{"the database is slow", "the database is slow"},
		{"dbx and nilpotent", "dbx and nilpotent"},
		{"err=nil", "err=null"},
		{"auth check on auth_service", "authentication check on auth_service"},
		{"", ""},
	}
	for _, c := range cases {
		if got := replaceLexicon(c.in, lex); got != c.want {
			t.Errorf("replaceLexicon(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestReplaceLexiconEmpty(t *testing.T) {
	if got := replaceLexicon("hello", nil); got != "hello" {
		t.Errorf("empty lexicon changed text: %q", got)
	}
}
