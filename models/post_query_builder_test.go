package models

import (
	"testing"
)

func TestWeCanValidateOurPostQueryBuilder(t *testing.T) {
	urlParams := map[string][]string{
		"sort": []string{"invalid"},
	}
	qb := postQueryBuilder{}
	err := qb.validate(urlParams)
	if err == nil {
		t.Error("Passing an invalid option for the 'sort' key should have thrown an error")
	}

	// add more checks here for other params
}

// we should create additional tests for the other default sort values
func TestWeCanSetTheDefaultValueForSortingPosts(t *testing.T) {
	qb := postQueryBuilder{}
	if qb.sort != "" {
		t.Error("The default sorting option should be an empty string")
	}

	urlParams := map[string][]string{
		"sort": []string{""},
	}
	err := qb.validate(urlParams)
	if err != nil {
		t.Error("An empty string should be a valid value for the sort option")
	}
	if qb.sort != "popular" {
		t.Errorf("After successful validation the sort option should be set to 'popular'. Current value: %v", qb.sort)
	}
}
