package models

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"

	"cloud.google.com/go/firestore"
)

type ClientForFirestore struct {
	ProjectID      string
	CredentialFile string
}

func (p *ClientForFirestore) NewClient(ctx context.Context) (*firestore.Client, error) {
	if p.CredentialFile != "" {
		if f, err := os.Stat(p.CredentialFile); err == nil && !f.IsDir() {
			app, err := firestore.NewClient(ctx, p.ProjectID, option.WithCredentialsFile(p.CredentialFile))
			return app, err
		}
	}

	app, err := firestore.NewClient(ctx, p.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return app, err
}

// Get dataはpointer, 参照渡し
func (p *ClientForFirestore) Get(ctx context.Context, colName, docKey string, data any) error {
	client, err := p.NewClient(ctx)
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
func (p *ClientForFirestore) Set(ctx context.Context, colName, docKey string, data any) error {
	client, err := p.NewClient(ctx)
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

func (p *ClientForFirestore) IsExist(ctx context.Context, colName string, docKeys ...string) (isExistKeys []string, err error) {
	client, err := p.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firestore: %v", err)
	}
	defer client.Close()

	for _, key := range docKeys {
		if _, err := client.Collection(colName).Doc(key).Get(ctx); err != nil {
			log.Debug().Str("function", "CheckExistKeysFirestore").Msgf("key: %s is ok, not exist", key)
			continue
		}

		// すでに存在するkeyを返却
		isExistKeys = append(isExistKeys, key)
	}

	return isExistKeys, nil
}
