package models

import (
	"context"
	"fmt"
	"os"
	"reflect"

	firebase "firebase.google.com/go"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

type ClientForFirebase struct {
	ProjectID      string
	CredentialFile string
}

func (p *ClientForFirebase) NewApp(ctx context.Context) (*firebase.App, error) {
	config := &firebase.Config{ProjectID: p.ProjectID}
	if p.CredentialFile != "" {
		if f, err := os.Stat(p.CredentialFile); err == nil && !f.IsDir() {
			app, err := firebase.NewApp(ctx, config, option.WithCredentialsFile(p.CredentialFile))
			return app, err
		}
	}

	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return app, err
}

// Get dataはpointer, 参照渡し
func (p *ClientForFirebase) Get(colName, docKey string, data any) error {
	ctx := context.Background()
	app, err := p.NewApp(ctx)
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

// Set dataはnot pointer, 値渡し
func (p *ClientForFirebase) Set(colName, docKey string, data any) error {
	ctx := context.Background()
	app, err := p.NewApp(ctx)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("error initializing firestore: %v", err)
	}
	defer client.Close()

	switch value := data.(type) {
	case []Account:
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

	case []Post:
		// 顧客投稿データの登録
		// id, keyはuuidで生成
		// 現状、投稿分の重複は容認、考慮せず
		for _, v := range value {
			v.UUID = uuid.New().String()
			v.SetCreateAt()

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

func (p *ClientForFirebase) IsExist(colName string, docKeys []string) error {
	ctx := context.Background()
	app, err := p.NewApp(ctx)
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
