package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/go-numb/go-spread-utils/cloud_run/post/models"

	zerolog "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/joho/godotenv/autoload"
)

const (
	IS_PRODUCTION = false
)

var (
	PORT string
)

func init() {
	if IS_PRODUCTION {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// 環境別の処理
	if runtime.GOOS == "linux" {
		PORT = fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
		log.Debug().Msgf("Linuxでの処理, PORT: %s", PORT)
	} else {
		PORT = fmt.Sprintf("localhost:%s", os.Getenv("PORT"))
		log.Debug().Msgf("その他のOSでの処理, PORT: %s", PORT)
	}

	// Twitter env for Twitter client@gotwi
	os.Setenv("GOTWI_API_KEY", os.Getenv("CONSUMERKEY"))
	os.Setenv("GOTWI_API_KEY_SECRET", os.Getenv("CONSUMERSECRET"))
}

func main() {
	http.HandleFunc("/post", handler)

	// Start HTTP server.
	log.Debug().Msgf("listening on port %s", PORT)
	log.Fatal().Err(http.ListenAndServe(PORT, nil)).Msg("http server error")
}

func handler(w http.ResponseWriter, r *http.Request) {
	// get query

	// get time
	t := time.Now()
	if err := Post(t); err != nil {
		log.Err(err).Str("function", "Post").TimeDiff("elapstime", time.Now(), t).Msg("post error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
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
	accounts, err := c.GetAccounts(t.Hour())
	if err != nil {
		return err
	}

	// Select accounts
	// 条件に合うアカウントを選択
	accounts, err = models.SelectAccount(accounts, t)
	if err != nil {
		return err
	}

	// 各ユーザーの処理
	for i := 0; i < len(accounts); i++ {
		// Get account spreadsheet:sheet@posts
		// 各ユーザーのスプレッドシート:postsを取得
		posts, err := c.GetPosts()
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

	return nil
}
