package models

// Vote is a positive (upvote) or negative (downvote) element that users can
// create for various other models in the system
// The Type column will reference what is being voted upon
type Vote struct {
	ID      uint   `db:"id"`
	VoterID uint   `db:"voter_id"`
	Vote    int    `db:"vote"`
	ItemID  uint   `db:"item_id"`
	Type    string `db:"type"`
}
