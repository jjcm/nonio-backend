package models

import (
	"errors"
	"net/url"
)

// postQueryBuilder is an object that we can set up to generate the correct SQL
// query when querying the DB for Post objects
type postQueryBuilder struct {
	sort     string
	offset   int
	time     string
	tag      string
	userName string
}

func (qb *postQueryBuilder) selectSQL(queryParams url.Values) (string, []interface{}) {
	var args []interface{}
	query := "SELECT * FROM `posts`"

	// set sort param
	query += " ORDER BY " + queryParams.Get("sort")

	// generate more SQL here to return the correct select statement
	//

	return query, args
}

func (qb *postQueryBuilder) validate(queryParams url.Values) error {
	validSortOptions := []string{
		"",
		"popular",
		"top",
		"new",
	}
	qb.sort = queryParams.Get("sort")
	if !stringInSlice(qb.sort, validSortOptions) {
		return errors.New("Invalid option passed for 'sort'")
	}
	// set default option for sort
	if qb.sort == "" {
		qb.sort = "popular"
	}

	// put more validation checks here, and set the appropriate defaults

	return nil
}
