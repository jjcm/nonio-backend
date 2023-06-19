package models

import (
	"fmt"
	"time"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/transfer"
)

type Payout struct {
	StripeConnectAccountId string
	Payout                 float64
}

func AllocatePayouts() error {
	currentTime := time.Now()

	ledgerEntries, err := calculatePayouts(currentTime)
	if err != nil {
		Log.Errorf("Error calculating payouts: %v", err)
		return err
	}

	for _, ledger := range ledgerEntries {
		tx, txErr := DBConn.Begin()
		if txErr != nil {
			return txErr
		}
		_, txErr = tx.Exec("insert into ledger (author_id, contributor_id, amount, type, description) values (?, ?, ?, ?, ?)",
			ledger.authorId, ledger.contributorId, ledger.amount, ledger.ledgerType, ledger.description,
		)
		_, txErr = tx.Exec("UPDATE users SET cash = cash + ? WHERE id = ?", ledger.amount, ledger.authorId)
		txErr = tx.Commit()
		if txErr != nil {
			return txErr
		}
	}

	u := User{}
	allUsers, err := u.GetAllForPayout()
	for _, user := range allUsers {
		tempCash := user.Cash
		params := &stripe.TransferParams{
			Amount:      stripe.Int64(int64(tempCash)),
			Currency:    stripe.String(string(stripe.CurrencyUSD)),
			Destination: stripe.String(user.StripeConnectAccountID),
		}

		_, err := transfer.New(params)
		if err != nil {
			return err
		}

		tx, txErr := DBConn.Begin()
		if txErr != nil {
			return txErr
		}
		_, txErr = tx.Exec("insert into ledger (author_id, contributor_id, amount, type, description) values (?, ?, ?, ?, ?)",
			user.ID, -1, tempCash, "withdrawal", "withdrawal from non.io to Stripe",
		)
		_, txErr = tx.Exec("UPDATE users SET cash = 0 WHERE id = ?", user.ID)
		txErr = tx.Commit()
		if txErr != nil {
			return txErr
		}

		err = user.UpdateLastPayout(time.Now())
		if err != nil {
			return err
		}

		now := time.Now()
		tomorrow := now.Add(time.Hour * 24)
		if user.CurrentPeriodEnd.After(now) && user.CurrentPeriodEnd.Before(tomorrow) {
			err = user.UpdateNextPayout(user.CurrentPeriodEnd)
			if err != nil {
				return err
			}
		}
	}

	v := PostTagVote{}
	err = v.MarkVotesAsTallied(currentTime)
	if err != nil {
		return err
	}

	return nil
}

func calculatePayouts(currentTime time.Time) ([]LedgerEntries, error) {
	u := User{}
	users, err := u.GetAllForPayout()
	if err != nil {
		Log.Error("Error getting users for payout")
		return nil, err
	}

	var ledgerEntries []LedgerEntries

	for _, user := range users {
		if user.AccountType == "supporter" {
			votes, err := user.GetUntalliedVotes(currentTime)
			if err != nil {
				Log.Error("Error getting users votes")
				return nil, err
			}
			votes = uniquePostFilter(votes)
			payoutPerVote := (user.SubscriptionAmount - ServerFee) / float64(len(votes))

			for _, vote := range votes {
				post := Post{}
				post.FindByID(vote.PostID)
				u.FindByID(post.AuthorID)

				ledgerEntries = append(ledgerEntries, LedgerEntries{
					authorId:      u.ID,
					contributorId: user.ID,
					amount:        payoutPerVote,
					ledgerType:    "deposit",
					description:   "deposit from user ID " + fmt.Sprintf("%d", user.ID),
				})
			}
		}
	}

	return ledgerEntries, nil
}

func uniquePostFilter(votes []PostTagVote) []PostTagVote {
	keys := make(map[int]bool)
	uniqueVotes := []PostTagVote{}

	for _, vote := range votes {
		if _, added := keys[vote.PostID]; !added {
			keys[vote.PostID] = true
			uniqueVotes = append(uniqueVotes, vote)
		}
	}

	return uniqueVotes
}

type LedgerEntries struct {
	authorId      int
	contributorId int
	amount        float64
	ledgerType    string
	description   string
}
