package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-numb/gcloud-spread-tweets/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// ScheduleProcess handles the form submission
// 1. 現在時刻を取得
// 2. 指定時間が合致するアカウントを取得
// 3. アカウントに紐づく投稿データを取得
// 4. 投稿データを選別
// 5. 投稿処理をリクエスト（別インスタンス
// 6. 投稿データを更新
func (p *Client) ScheduleProcess(c echo.Context) error {
	now := time.Now()
	ctx := context.Background()

	log.Debug().Msgf("schedule process: %s", now.String())

	// Get accounts
	accounts, err := p.getAccounts(ctx, now)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to get accounts: %v", err))
	}

	// Select accounts

	for _, account := range accounts {
		// Get posts
		posts, err := p.getPosts(ctx, account)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to get posts: %v", err))
		}

		// Select post
		post := posts[0]

		// Request to post process
		if err := p.RequestForPost(post.ToURLValues()); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to request to post process: %v", err))
		}

		// Update post database
		post.SetLastPostedAt()
		if err := p.Firestore.Set(ctx, Posts, post.UUID, post); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to update post: %v", err))
		}
	}

	return nil
}

// getAccounts カスタム検索
func (p *Client) getAccounts(ctx context.Context, t time.Time) ([]models.Account, error) {
	store, err := p.Firestore.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer store.Close()

	hour := t.Hour()
	minute := t.Minute()

	docs, err := store.Collection(Accounts).
		Where("hours", "array-contains", hour).
		Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var accounts []models.Account
	for _, doc := range docs {
		var account models.Account
		if err := doc.DataTo(&account); err != nil {
			log.Warn().Err(err).Msg("failed to convert data to account")
			continue
		}

		for _, m := range account.Minutes {
			if m != minute {
				continue
			}
			accounts = append(accounts, account)
		}
	}

	return accounts, nil
}

// getPosts カスタム検索
func (p *Client) getPosts(ctx context.Context, account models.Account) ([]models.Post, error) {
	store, err := p.Firestore.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer store.Close()

	docs, err := store.Collection(Posts).
		Where("id", "==", account.ID).
		// Where("checked", "!=", 0). // Checked=0でも取得（GUIで確認変更を可能にするため）
		Where("is_delete", "==", false).
		Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var posts []models.Post
	for _, doc := range docs {
		var post models.Post
		if err := doc.DataTo(&post); err != nil {
			log.Warn().Err(err).Msg("failed to convert data to posts")
			continue
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// isOK Token&Usernameが存在するか確認
func (p *Client) isOK(ctx context.Context, token, username string) bool {
	store, err := p.Firestore.NewClient(ctx)
	if err != nil {
		return false
	}
	defer store.Close()

	doc, err := store.Collection(Users).Doc(token).Get(ctx)
	if err != nil {
		return false
	}

	var claims models.Claims
	if err := doc.DataTo(&claims); err != nil {
		log.Warn().Err(err).Msg("failed to convert data to posts")
		return false
	}

	if claims.ID != username {
		return false
	}

	return true
}
