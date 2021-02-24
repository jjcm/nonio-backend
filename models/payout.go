package models

import (
	"fmt"
	"time"
)

type Payout struct {
	UserID int
	Payout float64
}

func CalculatePayouts() ([]Payout, error) {
	currentTime := time.Now()
	fmt.Printf("Routine ran at %v\n", currentTime.String())

	u := User{}
	users, err := u.GetAll()
	var payouts []Payout
	if err != nil {
		Log.Error("Error getting list of users")
		return payouts, err
	}

	// For each of our users, get their votes and calculate what their individual payout is.
	for _, user := range users {
		votes, err := user.GetUntalliedVotes(currentTime)
		fmt.Printf("Server fee is %v, %v's subscription is %v and has %v votes.\n", ServerFee, user.Username, user.SubscriptionAmount, len(votes))
		payoutPerVote := (user.SubscriptionAmount - ServerFee) / float64(len(votes))
		fmt.Printf("Payout per vote for user %v is %v\n", user.Username, payoutPerVote)
		if err != nil {
			Log.Errorf("Error getting votes for user %v\n", user.Email)
			return payouts, err
		}
		for _, vote := range votes {
			fmt.Println(vote.PostID)
		}
	}
	return payouts, err
}
