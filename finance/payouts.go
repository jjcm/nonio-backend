package finance

import (
	"fmt"
	"soci-backend/models"
	"time"
)

type Payout struct {
	UserID int
	Payout float64
}

func CalculatePayouts() ([]Payout, error) {
	currentTime := time.Now()
	fmt.Printf("Routine ran at %v\n", currentTime.String())

	u := models.User{}
	users, err := u.GetAll()
	var payouts []Payout
	if err != nil {
		Log.Error("Error getting list of users")
		return payouts, err
	}

	// For each of our users, get their votes and calculate what their individual payout is.
	for _, user := range users {
		votes, err := user.GetUntalliedVotes(currentTime)
		fmt.Println(ServerFee)
		payoutPerVote := (user.SubscriptionAmount - ServerFee) / float64(len(votes))
		fmt.Println(payoutPerVote)
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

// DemoInsertUser inserts an example user into the db
func DemoInsertUser() error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := DBConn.Exec("INSERT INTO users (email, username, password, subscription_amount, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)", "example@example.com", "example", "asdf", 10, now, now)
	if err != nil {
		return err
	}
	return nil
}
