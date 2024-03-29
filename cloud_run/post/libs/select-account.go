package libs

import (
	"fmt"
	"time"

	"github.com/go-numb/gcloud-spread-tweets/models"
)

var EventKeys = []string{"all", "-", ""}

func SelectAccount(accounts []models.Account, t time.Time) ([]models.Account, error) {
	var results []models.Account
	for i := 0; i < len(accounts); i++ {
		if accounts[i].ID == "" {
			continue
		}

		if accounts[i].AccessToken == "" {
			continue
		}

		if accounts[i].AccessSecret == "" {
			continue
		}

		if accounts[i].Subscribed == 0 {
			continue
		}

		if !isTargetMinutes(accounts[i].Minutes, t.Minute()) {
			continue
		}

		results = append(results, accounts[i])
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("accounts is empty")
	}

	return results, nil
}

func isTargetMinutes(targetMinute []int, nowMinute int) bool {
	for i := 0; i < len(targetMinute); i++ {
		if targetMinute[i] == nowMinute {
			return true
		}
	}

	return false
}
