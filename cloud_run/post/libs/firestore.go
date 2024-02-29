package libs

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (p *Client) GetAccounts(hour int) ([]Account, error) {
	ctx := context.Background()
	docs, err := p.Firestore.Collection("accounts").Where("hours", "array-contains", hour).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var accounts = make([]Account, len(docs))
	for _, doc := range docs {
		var account Account
		if err := doc.DataTo(&account); err != nil {
			log.Warn().Err(err).Msg("failed to convert data to account")
			continue
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (p *Client) GetPosts() ([]Post, error) {
	ctx := context.Background()
	docs, err := p.Firestore.Collection("posts").Where("checked", "!=", 0).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var posts = make([]Post, len(docs))
	for _, doc := range docs {
		var post Post
		if err := doc.DataTo(&post); err != nil {
			log.Warn().Err(err).Msg("failed to convert data to posts")
			continue
		}

		posts = append(posts, post)
	}

	return posts, nil
}
