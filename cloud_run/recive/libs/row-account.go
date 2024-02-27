package libs

import "github.com/google/uuid"

type Account struct {
	UUID        string `csv:"uuid" dataframe:"uuid"`
	ID          string `csv:"x_id" dataframe:"x_id"`
	Password    string `csv:"x_password" dataframe:"x_password"`
	SpreadID    string `csv:"spread_id" dataframe:"spread_id"`
	AccessToken string `csv:"access_token" dataframe:"access_token"`
	SecretToken string `csv:"secret_token" dataframe:"secret_token"`
	Subscribed  int    `csv:"subscribed" dataframe:"subscribed"`
	Hours       string `csv:"hours" dataframe:"hours"`
	Minutes     string `csv:"minutes" dataframe:"minutes"`
	TermHours   int    `csv:"term_hours" dataframe:"term_hours"`
}

// NewAccountForFree 初回会員登録用
func NewAccountForFree(id, sheetID, accessToken, accessSecret string) *Account {
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

// NewAccountForSubscribe 課金顧客用
// TODO: 何を提供するか未定
func NewAccountForSubscribe(id, sheetID, accessToken, accessSecret string) *Account {
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
		TermHours:   24,
	}
}

// NewAccountForPro 高課金顧客用
// TODO: 何を提供するか未定
func NewAccountForPro(id, sheetID, accessToken, accessSecret string) *Account {
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
		TermHours:   24,
	}
}
