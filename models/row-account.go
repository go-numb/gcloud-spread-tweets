package models

import (
	"github.com/google/uuid"
)

const (
	SubscribedFree SubscribedPlan = iota
	SubscribedBasic
	SubscribedPro
)

type SubscribedPlan uint8

type Account struct {
	UUID         string         `csv:"uuid" dataframe:"uuid" firestore:"uuid,omitempty" json:"uuid,omitempty"`
	ID           string         `csv:"id" dataframe:"id" firestore:"id,omitempty" json:"id,omitempty"`
	Password     string         `csv:"password" dataframe:"password" firestore:"password,omitempty" json:"password,omitempty"`
	SpreadID     string         `csv:"spread_id" dataframe:"spread_id" firestore:"spread_id,omitempty" json:"spread_id,omitempty"`
	AccessToken  string         `csv:"access_token" dataframe:"access_token" firestore:"access_token,omitempty" json:"access_token,omitempty"`
	AccessSecret string         `csv:"access_secret" dataframe:"access_secret" firestore:"access_secret,omitempty" json:"access_secret,omitempty"`
	Subscribed   SubscribedPlan `csv:"subscribed" dataframe:"subscribed" firestore:"subscribed,omitempty" json:"subscribed,omitempty"`
	Hours        []int          `csv:"hours" dataframe:"hours" firestore:"hours,omitempty" json:"hours,omitempty"`
	Minutes      []int          `csv:"minutes" dataframe:"minutes" firestore:"minutes,omitempty" json:"minutes,omitempty"`
	TermHours    int            `csv:"term_hours" dataframe:"term_hours" firestore:"term_hours,omitempty" json:"term_hours,omitempty"`
}

// NewAccountForFree 初回会員登録用
func NewAccount(id, sheetID, accessToken, accessSecret string) *Account {
	return &Account{
		UUID:         uuid.New().String(),
		ID:           id,
		Password:     "",
		SpreadID:     sheetID,
		AccessToken:  accessToken,
		AccessSecret: accessSecret,
		Subscribed:   0,
		Hours:        make([]int, 0),
		Minutes:      make([]int, 0),
		TermHours:    48,
	}
}

// GetID for interface
func (p Account) GetID() string {
	return p.ID
}

// SetSubscribed 有料プランの設定
func (a *Account) SetSubscribed(level SubscribedPlan) *Account {
	a.Subscribed = level
	return a
}

// SetTime 投稿時間の設定
func (a *Account) SetTime(hours, minutes []int) *Account {
	if len(hours) != 0 {
		a.Hours = hours
	}
	if len(minutes) != 0 {
		a.Minutes = minutes
	}
	return a
}

// GetTime 投稿時間の取得
func (a *Account) GetTime() (hours, minutes []int) {
	return a.Hours, a.Minutes
}

// SetTerm 再投稿間隔設定
func (a *Account) SetTerm(hours int) *Account {
	a.TermHours = hours
	return a
}

// func intSliceToString(numbers []int) string {
// 	str := make([]string, len(numbers))
// 	for i := 0; i < len(numbers); i++ {
// 		str[i] = strconv.Itoa(numbers[i])
// 	}
// 	return strings.Join(str, ",")
// }

// func stringsToIntSlice(s string) []int {
// 	str := strings.Split(s, ",")
// 	num := make([]int, len(str))
// 	for i := 0; i < len(str); i++ {
// 		n, err := strconv.Atoi(strings.TrimSpace(str[i]))
// 		if err != nil {
// 			continue
// 		}
// 		num[i] = n
// 	}
// 	return num
// }
