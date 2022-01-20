package models

import (
	"fmt"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/transfer"
	"time"
)

type Payout struct {
	StripeConnectAccountId string
	Payout                 float64
}

func AllocatePayouts() error {
	currentTime := time.Now()
	payouts, err := calculatePayouts(currentTime)
	if err != nil {
		Log.Errorf("Error calculating payouts: %v", err)
		return err
	}

	for userId, payout := range payouts {
		_, err := DBConn.Exec("UPDATE users SET cash = cash + ? WHERE id = ?", payout, userId)
		if err != nil {
			return err
		}
	}

	u := User{}
	allUsers, err := u.GetAllForPayout()
	if err != nil {
		return err
	}

	for _, user := range allUsers {
		fmt.Printf("User %v is getting paid %v\n", user.ID, user.Cash)

		params := &stripe.TransferParams{
			Amount:      stripe.Int64(int64(user.Cash)),
			Currency:    stripe.String(string(stripe.CurrencyUSD)),
			Destination: stripe.String(user.StripeConnectAccountID),
		}

		_, err := transfer.New(params)
		if err != nil {
			return err
		}
		_, err = DBConn.Exec("UPDATE users SET cash = 0 WHERE id = ?", user.ID)
		if err != nil {
			return err
		}
	}

	v := PostTagVote{}
	err = v.MarkVotesAsTallied(currentTime)
	if err != nil {
		return err
	}

	return nil
}

func calculatePayouts(currentTime time.Time) (map[int]float64, error) {
	fmt.Printf("Routine ran at %v\n", currentTime.String())

	u := User{}
	users, err := u.GetAll()
	payouts := map[int]float64{}
	if err != nil {
		Log.Errorf("Error getting list of users: %v\n", err)
		return nil, err
	}

	// For each of our users, get their votes and calculate what their individual payout is.
	for _, user := range users {
		if user.AccountType == "supporter" {
			votes, err := user.GetUntalliedVotes(currentTime)

			// A user may have multiple tags they voted on for a post. A vote for a post should only be counted once, regardless of the tags upvoted.
			votes = uniquePostFilter(votes)

			payoutPerVote := (user.SubscriptionAmount - ServerFee) / float64(len(votes))

			if err != nil {
				Log.Errorf("Error getting votes for user %v\n", user.Email)
				return nil, err
			}

			for _, vote := range votes {
				post := Post{}
				post.FindByID(vote.PostID)
				u.FindByID(post.AuthorID)

				payouts[u.ID] += payoutPerVote
			}
		}
	}
	return payouts, err
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
