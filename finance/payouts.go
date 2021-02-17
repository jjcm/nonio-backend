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
	fmt.Println("created user")
	users, err := u.GetAll()
	fmt.Println("assigned users")
	if err != nil {
		fmt.Println("error with user")
		fmt.Println(err)
		return
	}

	fmt.Println("printing users")
	for _, user := range users {
		fmt.Println(user.Email)
	}
}
