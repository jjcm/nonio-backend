package models

import "testing"

func TestWeCanCreateAWebFriendlyAliasFromAGivenString(t *testing.T) {
	cases := map[string]string{
		"Hello":                      "hello",
		"This is a TEST":             "this-is-a-test",
		"   extra spaces are NBD   ": "extra-spaces-are-nbd",
		"spaces  between  words":     "spaces-between-words",
		"ditch bad chars !@#$%^&*()_+-=,./<>?;'\"[]{}`~": "ditch-bad-chars-_-.~",
	}

	for title, alias := range cases {
		if createURLFromTitle(title) != alias {
			t.Errorf("Expected alias didn't match our title.\nTitle: '%v'\nExpected Alias:  %v\nGenerated Alias: %v", title, alias, createURLFromTitle(title))
		}
	}
}
