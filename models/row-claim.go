package models

type Claims struct {
	Name string `json:"name"`

	AccessToken  string `json:"access_token"`
	AccessSecret string `json:"access_secret"`

	// Auth Request Token
	RequestToken       string `json:"request_token"`
	RequestTokenSecret string `json:"request_token_secret"`

	SpreadsheetID string `json:"spreadsheet_id"`
}

// GetID for interface
func (p Claims) GetID() string {
	return p.Name
}
