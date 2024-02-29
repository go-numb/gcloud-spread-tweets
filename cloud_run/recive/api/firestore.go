package api

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/go-numb/gcloud-spread-tweets/models"

	firebase "firebase.google.com/go/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

const (
	// firestore collection
	Users    = "users"
	Accounts = "accounts"
	Posts    = "posts"

	// bind key for map
	RTOKEN       = "request_token"
	RTOKENSECRET = "request_token_secret"

	ACCESSTOKEN       = "access_token"
	ACCESSTOKENSECRET = "access_token_secret"
)

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

	switch value := data.(type) {
	case []models.Account:
		// 顧客アカウントの登録
		// id, keyはTwitter/Xアカウントのidで生成
		// 現状、アカウント分の重複は容認せず
		for _, v := range value {
			// Twitter/X IDをキーにして重複を許さない
			if _, err := client.Collection(colName).Doc(v.ID).Set(ctx, v); err != nil {
				log.Error().Err(err).Msgf("error setting document: data type %s", reflect.TypeOf(v).String())
				continue
			}
		}

	case []models.Post:
		// 顧客投稿データの登録
		// id, keyはuuidで生成
		// 現状、投稿分の重複は容認、考慮せず
		for _, v := range value {
			v.UUID = uuid.New().String()
			if _, err := client.Collection(colName).Doc(v.UUID).Set(ctx, v); err != nil {
				log.Error().Err(err).Msgf("error setting document: data type %s", reflect.TypeOf(v).String())
				continue
			}
		}

	default:
		if _, err := client.Collection(colName).Doc(docKey).Set(ctx, value); err != nil {
			return fmt.Errorf("error setting document: %v, data type: %s", err, reflect.TypeOf(data).String())
		}
	}

	return nil
}

func (p *Client) CheckExistKeysFirestore(colName string, docKeys []string) error {
	app, ctx, err := p.NewAppClient()
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("error initializing firestore: %v", err)
	}
	defer client.Close()

	for _, key := range docKeys {
		if _, err := client.Collection(colName).Doc(key).Get(ctx); err != nil {
			// log.Debug().Str("function", "CheckExistKeysFirestore").Msgf("key: %s is ok, not exist", key)
			continue
		}

		// 一つでもアカウント重複があればエラー
		return fmt.Errorf("error checking exist keys: %s", key)
	}

	return nil
}
