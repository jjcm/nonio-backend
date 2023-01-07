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


	payouts, ledgerEntries, err := calculatePayouts(currentTime)
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

		err = user.UpdateLastPayout(time.Now())
		if err != nil {
			return err
		}

		for _, l := range ledgerEntries {
			if l.authorId == user.ID {
				// deposit ledgers
				return createLedgerEntry(user.ID, l.contributorId, user.Cash, l.ledgerType, l.description)
			}
		}

		//withdrawal ledger entry
		err = createLedgerEntry(user.ID, -1, user.Cash, "withdrawal", "withdrawal from non.io to Stripe")
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

func calculatePayouts(currentTime time.Time) (map[int]float64, []LedgerEntries, error) {
	fmt.Printf("Routine ran at %v\n", currentTime.String())

	u := User{}
	users, err := u.GetAllForPayout()
	payouts := map[int]float64{}
	if err != nil {
		Log.Errorf("Error getting list of users: %v\n", err)
		return nil, nil, err
	}

	// For each of our users, get their votes and calculate what their individual payout is.
	var ledgerEntries []LedgerEntries
	for _, user := range users {
		if user.AccountType == "supporter" {
			votes, err := user.GetUntalliedVotes(currentTime)

			// A user may have multiple tags they voted on for a post. A vote for a post should only be counted once, regardless of the tags upvoted.
			votes = uniquePostFilter(votes)

			payoutPerVote := (user.SubscriptionAmount - ServerFee) / float64(len(votes))

			if err != nil {
				Log.Errorf("Error getting votes for user %v\n", user.Email)
				return nil, nil, err
			}

			for _, vote := range votes {
				post := Post{}
				post.FindByID(vote.PostID)
				u.FindByID(post.AuthorID)

				payouts[u.ID] += payoutPerVote
			}

			ledgerEntries = append(ledgerEntries, LedgerEntries{
				authorId:      u.ID,
				contributorId: user.ID,
				amount:        payouts[u.ID],
				ledgerType:    "deposit",
				description:   "deposit from " + user.Name,
			})
		}
	}
	return payouts, ledgerEntries, err
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

func createLedgerEntry(authorId, contributorId int, amount float64, ledgerType, description string) error {
	_, err := DBConn.Exec("insert into ledger (author_id, contributor_id, amount, type, description) values (?, ?, ?, ?, ?)",
		authorId, contributorId, amount, ledgerType, description,
	)
	if err != nil {
		return err
	}

	return nil
}

/* This is proposed pseudocode for this. Lemme know if you think this would work/make sense */
AllocatePayouts () {
	// we only grab the ledger entries for the current period. See the modified calculatePayouts function below.
	ledgerEntries, err := calculatePayouts(currentTime)

	// using pseudocode here just to communicate the idea
	for each ledger in ledgerEntries {

		//create a transaction that creates a ledger, and increments the user's cash
		transaction {
			db(insert into ledgertable deposit ${ledger})
			db(update users set cash = cash + ${ledger.amount} where id = ${ledger.authorId})
		}

	}
	
	// after this we can run the loop that pays out users to stripe if they have that set up

	allUsers, err := u.GetAllForPayout()
	for each user in allUsers {
		
		// store some temporary cash value
		tempCash = user.cash

		// run a transaction that creates a ledger, and sets the user's cash to 0
		transaction {
			db(insert into ledgertable withdrawal of ${user.cash})
			db(update users set cash = 0 where id = ${user.id})
		}

		// if no errors, then transfer our tempCash to the user's stripe account
		if(!err){
			stripe.transfer(tempCash, user.stripeConnectId)
		}


	}

}

func calculatePayouts(currentTime time.Time) ([]LedgerEntries, error) {
	users, err := u.GetAllForPayout()
	var ledgerEntries []LedgerEntries

	// For each of our users, get their votes and calculate what their individual payout is.
	for _, user := range users {
		if user.AccountType == "supporter" {

			// get all of our untallied votes for the user
			votes, err := user.GetUntalliedVotes(currentTime)

			// A user may have multiple tags they voted on for a post. A vote for a post should only be counted once, regardless of the tags upvoted.
			votes = uniquePostFilter(votes)

			// This is key - we still calculate the payout per vote. 
			payoutPerVote := (user.SubscriptionAmount - ServerFee) / float64(len(votes))

			// then for each unique vote, we create a ledger entry
			for _, vote := range votes {
				// find the author of the post
				post := Post{}
				post.FindByID(vote.PostID)
				u.FindByID(post.AuthorID)

				// now create a ledger entry with all of our info. The amount will be the "payoutPerVote" amount we calculated above.
				ledgerEntries = append(ledgerEntries, LedgerEntries{
					authorId:      u.ID,
					contributorId: user.ID,
					amount:        payoutPerVote,
					ledgerType:    "deposit",
					description:   "deposit from " + user.Name,
				})
			}

		}
	}
	return ledgerEntries, err
}