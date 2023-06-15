package models

// URLIsAvailable - check the database to see if a given URL is available
func URLIsAvailable(url string) (bool, error) {
	// remove the leading "/post/url-is-available/"
	url = url[23:]

	var total int
	err := DBConn.Get(&total, "SELECT COUNT(*) FROM posts where url = ?", url)
	if err != nil {
		return false, err
	}
	if total != 0 {
		return false, nil
	}
	return true, nil
}
