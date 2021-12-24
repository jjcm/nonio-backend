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

	// map of stripe connect account id and the corresponding payout to be paid out
	users := make(map[string]float64)

	for _, payout := range payouts {
		users[payout.StripeConnectAccountId] = users[payout.StripeConnectAccountId] + payout.Payout
	}

	for stripeAccountId, payout := range users {
		fmt.Printf("User %v is getting paid %v\n", stripeAccountId, payout)

		params := &stripe.TransferParams{
			Amount:      stripe.Int64(int64(payout)),
			Currency:    stripe.String(string(stripe.CurrencyUSD)),
			Destination: stripe.String(stripeAccountId),
		}

		tr, _ := transfer.New(params)

		fmt.Println(tr.ID)
	}

	v := PostTagVote{}
	err = v.MarkVotesAsTallied(currentTime)
	if err != nil {
		return err
	}

	return nil
}

func calculatePayouts(currentTime time.Time) ([]Payout, error) {
	fmt.Printf("Routine ran at %v\n", currentTime.String())

	u := User{}
	users, err := u.GetAll()
	var payouts []Payout
	if err != nil {
		Log.Errorf("Error getting list of users: %v\n", err)
		return payouts, err
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
				return payouts, err
			}

			for _, vote := range votes {
				post := Post{}
				post.FindByID(vote.PostID)
				u.FindByID(post.AuthorID)
				payout := Payout{
					StripeConnectAccountId: u.StripeConnectAccountID,
					Payout:                 payoutPerVote,
				}

				payouts = append(payouts, payout)
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
