package models

import (
	"encoding/json"
	"time"
)

// Tag - code representation of a single tag
type Ledger struct {
	ID            int       `db:"id" json:"-"`
	AuthorID      int       `db:"author_id" json:"-"`
	ContributorID int       `db:"contributor_id" json:"-"`
	Type          string    `db:"type" json:"type"`
	Amount        float64   `db:"amount" json:"amount"`
	Description   string    `db:"description" json:"description"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
}

// ToJSON - get a string representation of this Tag in JSON
func (t *Tag) ToJSON() string {
	jsonData, err := json.Marshal(t)
	if err != nil {
		return err.Error()
	}
	return string(jsonData)
}

/************************************************/
/******************** CREATE ********************/
/************************************************/

// CreateLedger - create a new ledger entry in the database
func CreateLedger(author User, contributor User, ledgerType string, amount float64, description string) error {
	_, err := DBConn.Exec("INSERT INTO ledger (author_id, contributor_id, type, amount, description) VALUES (?, ?, ?, ?, ?)", author.ID, contributor.ID, ledgerType, amount, description)
	if err != nil {
		return err
	}
	return nil
}

func (l *Ledger) CreateLedgerWithTx(tx Transaction) error {
	_, err := tx.Exec("INSERT INTO ledger (author_id, contributor_id, type, amount, description) VALUES (?, ?, ?, ?, ?)", l.AuthorID, l.ContributorID, l.Type, l.Amount, l.Description)
	if err != nil {
		return err
	}

	return nil
}

/************************************************/
/********************* READ *********************/
/************************************************/

// GetLedgerEntries - get ledger entries for a specific user
func (u *User) GetLedgerEntries(name string) ([]Ledger, error) {
	ledgerEntries := []Ledger{}
	err := DBConn.Get(&ledgerEntries, "SELECT * FROM ledger WHERE author_id = ?", u.ID)
	if err != nil {
		return nil, err
	}

	return ledgerEntries, nil
}
