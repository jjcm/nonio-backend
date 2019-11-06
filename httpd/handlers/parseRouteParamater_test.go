package handlers

import "testing"

// TestWeCanParseAURLParamater here we will test our URL param parser
func TestWeCanParseAURLParamater(t *testing.T) {
	testCases := []struct {
		fullURL  string
		pattern  string
		expected string
	}{
		{"/posts/yippie", "/posts/", "yippie"},
	}

	for _, tc := range testCases {
		output := parseRouteParamater(tc.fullURL, tc.pattern)
		if output != tc.expected {
			t.Errorf("URL param was not caught correctly. Expected %v got %v", tc.expected, output)
		}
	}
}
