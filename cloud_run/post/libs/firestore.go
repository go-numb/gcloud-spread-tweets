package libs

import (
	"context"

	"github.com/go-numb/gcloud-spread-tweets/models"
	"github.com/rs/zerolog/log"
)

func (p *Client) GetAccounts(hour int) ([]models.Account, error) {
	ctx := context.Background()
	store, err := p.Firestore.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer store.Close()

	docs, err := store.Collection("accounts").Where("hours", "array-contains", hour).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var accounts = make([]models.Account, len(docs))
	for _, doc := range docs {
		var account models.Account
		if err := doc.DataTo(&account); err != nil {
			log.Warn().Err(err).Msg("failed to convert data to account")
			continue
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (p *Client) GetPosts(account models.Account) ([]Post, error) {
	ctx := context.Background()
	store, err := p.Firestore.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer store.Close()

	docs, err := store.Collection("posts").
		Where("checked", "!=", 0).
		Where("is_deleted", "!=", true).
		Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var posts = make([]Post, len(docs))
	for _, doc := range docs {
		var post Post
		if err := doc.DataTo(post.Post); err != nil {
			log.Warn().Err(err).Msg("failed to convert data to posts")
			continue
		}

		post.AccessToken = account.AccessToken
		post.AccessSecret = account.AccessSecret
		post.Password = account.Password

		posts = append(posts, post)
	}

	return posts, nil
}
