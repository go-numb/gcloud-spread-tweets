package libs

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	SubscribedFree SubscribedPlan = iota
	SubscribedBasic
	SubscribedPro
)

type SubscribedPlan uint8

type Account struct {
	UUID        string         `csv:"uuid" dataframe:"uuid"`
	ID          string         `csv:"x_id" dataframe:"x_id"`
	Password    string         `csv:"x_password" dataframe:"x_password"`
	SpreadID    string         `csv:"spread_id" dataframe:"spread_id"`
	AccessToken string         `csv:"access_token" dataframe:"access_token"`
	SecretToken string         `csv:"secret_token" dataframe:"secret_token"`
	Subscribed  SubscribedPlan `csv:"subscribed" dataframe:"subscribed"`
	Hours       string         `csv:"hours" dataframe:"hours"`
	Minutes     string         `csv:"minutes" dataframe:"minutes"`
	TermHours   int            `csv:"term_hours" dataframe:"term_hours"`
}

// NewAccountForFree 初回会員登録用
func NewAccount(id, sheetID, accessToken, accessSecret string) *Account {
	return &Account{
		UUID:        uuid.New().String(),
		ID:          id,
		Password:    "",
		SpreadID:    sheetID,
		AccessToken: accessToken,
		SecretToken: accessSecret,
		Subscribed:  0,
		Hours:       "",
		Minutes:     "",
		TermHours:   48,
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
		a.Hours = intSliceToString(hours)
	}
	if len(minutes) != 0 {
		a.Minutes = intSliceToString(minutes)
	}

	return a
}

// GetTime 投稿時間の取得
func (a *Account) GetTime() (hours, minutes []int) {
	if a.Hours != "" {
		hours = stringsToIntSlice(a.Hours)
	}
	if a.Minutes != "" {
		minutes = stringsToIntSlice(a.Minutes)
	}
	return
}

// SetTerm 再投稿間隔設定
func (a *Account) SetTerm(hours int) *Account {
	a.TermHours = hours
	return a
}

func intSliceToString(numbers []int) string {
	str := make([]string, len(numbers))
	for i := 0; i < len(numbers); i++ {
		str[i] = strconv.Itoa(numbers[i])
	}
	return strings.Join(str, ",")
}

func stringsToIntSlice(s string) []int {
	str := strings.Split(s, ",")
	num := make([]int, len(str))
	for i := 0; i < len(str); i++ {
		n, err := strconv.Atoi(strings.TrimSpace(str[i]))
		if err != nil {
			continue
		}
		num[i] = n
	}
	return num
}