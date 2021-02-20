package finance

import (
	"fmt"
	"soci-backend/models"
	"time"
)

func CalculatePayouts() {
	currentTime := time.Now()
	fmt.Printf("Routine ran at %v\n", currentTime.String())

	u := models.User{}
	users, err := u.GetAll()
	if err != nil {
		Log.Error("Error getting list of users")
		return
	}

	// For each of our users, get their votes and calculate what their individual payout is.
	for _, user := range users {
		votes, err := user.GetUntalliedVotes(currentTime)
		fmt.Println(ServerFee)
		payoutPerVote := (user.SubscriptionAmount - ServerFee) / float64(len(votes))
		fmt.Println(payoutPerVote)
		if err != nil {
			Log.Errorf("Error getting votes for user %v\n", user.Email)
			return
		}
		for _, vote := range votes {
			fmt.Println(vote.PostID)
		}
	}
}
