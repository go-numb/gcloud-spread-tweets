package models

import "fmt"

var EventKeys = []string{"all", "-", ""}

func SelectAccounts(accounts []Account, eventType string) ([]Account, error) {
	// すべて返す条件
	for i := 0; i < len(EventKeys); i++ {
		if eventType == EventKeys[i] {
			return accounts, nil
		}
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("accounts is empty")
	}

	return accounts, nil
}
