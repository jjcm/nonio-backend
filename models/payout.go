package models

import (
	"fmt"
	"time"
)

type Payout struct {
	UserID int
	Payout float64
}

func AllocatePayouts() error {
	currentTime := time.Now()
	payouts, err := calculatePayouts(currentTime)
	if err != nil {
		Log.Errorf("Error calculating payouts: %v", err)
		return err
	}

	users := make(map[int]float64)

	for _, payout := range payouts {
		users[payout.UserID] = users[payout.UserID] + payout.Payout
	}

	for user, payout := range users {
		fmt.Printf("User %v is getting paid %v\n", user, payout)
		_, err := DBConn.Exec("UPDATE users SET cash = cash + ? WHERE id = ?", payout, user)
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

			payout := Payout{
				UserID: post.AuthorID,
				Payout: payoutPerVote,
			}

			payouts = append(payouts, payout)
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
