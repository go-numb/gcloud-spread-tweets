package api

import (
	"context"
	"fmt"
	"os"
	"reflect"

	firebase "firebase.google.com/go"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

const (
	// firestore collection
	Users = "users"

	// bind key for map
	RTOKEN       = "request_token"
	RTOKENSECRET = "request_token_secret"

	ACCESSTOKEN       = "access_token"
	ACCESSTOKENSECRET = "access_token_secret"
)

type Claims struct {
	Name string `json:"name"`

	AccessToken  string `json:"access_token"`
	AccessSecret string `json:"access_secret"`

	// Auth Request Token
	RequestToken       string `json:"request_token"`
	RequestTokenSecret string `json:"request_token_secret"`

	SpreadsheetID string `json:"spreadsheet_id"`
}

func (p *Client) NewAppClient() (*firebase.App, context.Context, error) {
	config := &firebase.Config{ProjectID: p.ProjectID}
	ctx := context.Background()
	if p.CredentialFile != "" {
		if f, err := os.Stat(p.CredentialFile); err == nil && !f.IsDir() {
			app, err := firebase.NewApp(ctx, config, option.WithCredentialsFile(p.CredentialFile))
			return app, ctx, err
		}

		app, err := firebase.NewApp(ctx, config)
		return app, ctx, err
	}

	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		return nil, ctx, fmt.Errorf("error initializing app: %v", err)
	}

	return app, ctx, err
}

// GetAnyFirestore dataはpointer, 参照渡し
func (p *Client) GetAnyFirestore(colName, docKey string, data any) error {
	app, ctx, err := p.NewAppClient()
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("error initializing firestore: %v", err)
	}
	defer client.Close()

	doc, err := client.Collection(colName).Doc(docKey).Get(ctx)
	if err != nil {
		return fmt.Errorf("error setting document: %v", err)
	}

	// data bind
	if err := doc.DataTo(data); err != nil {
		return fmt.Errorf("error getting data: %v", err)
	}

	log.Debug().Msgf("get firestore, %+v, data type: %s", data, reflect.TypeOf(data).String())

	return nil
}

// SetAnyFirestore dataはnot pointer, 値渡し
func (p *Client) SetAnyFirestore(colName, docKey string, data any) error {
	app, ctx, err := p.NewAppClient()
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("error initializing firestore: %v", err)
	}
	defer client.Close()

	if _, err := client.Collection(colName).Doc(docKey).Set(ctx, data); err != nil {
		return fmt.Errorf("error setting document: %v, data type: %s", err, reflect.TypeOf(data).String())
	}

	return nil
}
