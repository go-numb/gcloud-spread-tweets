package main

import (
	"context"
	"fmt"
	"os"
	"post/models"
	"time"

	zerolog "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/joho/godotenv/autoload"
)

const (
	IS_PRODUCTION = false
)

func init() {
	if IS_PRODUCTION {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Twitter env for Twitter client@gotwi
	os.Setenv("GOTWI_API_KEY", os.Getenv("CONSUMERKEY"))
	os.Setenv("GOTWI_API_KEY_SECRET", os.Getenv("CONSUMERSECRET"))
}

func main() {
	t := time.Now()
	if err := Post(t); err != nil {
		log.Err(err).Str("function", "Post").TimeDiff("elapstime", time.Now(), t).Msg("post error")
		return
	}
}

func Post(t time.Time) error {
	if models.IsTimeToPost(t) {
		return fmt.Errorf("time switcher is off")
	}

	ctx, cancel := context.WithCancel(context.Background())
	c := models.NewClient(ctx)
	defer cancel()

	// Read spreadsheet:sheet@users rows
	// 管理者のスプレッドシート:usersを読み込む
	var accounts []models.Account
	accountsDf, err := c.GetAccounts(&accounts)
	if err != nil {
		return err
	}

	// Select accounts
	// 条件に合うアカウントを選択
	accounts, err = models.SelectAccounts(accounts, "all")
	if err != nil {
		return err
	}

	// 各ユーザーの処理
	for i := 0; i < len(accounts); i++ {
		// Get account spreadsheet:sheet@posts
		// 各ユーザーのスプレッドシート:postsを取得
		posts, err := c.GetPosts(accounts[i])
		if err != nil {
			return err
		}

		// Select posts
		// ユーザーが指定条件に合う投稿を選択
		post, err := models.SelectPost(posts)
		if err != nil {
			return err
		}

		// Post to Twitter/X
		// Twitter/Xに投稿
		if err := post.Do(); err != nil {
			return err
		}

		// Update spreadsheet:sheet@posts

	}

	var (
		_ = accountsDf
	)

	return nil
}
