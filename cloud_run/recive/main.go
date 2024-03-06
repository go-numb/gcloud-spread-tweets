package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-numb/gcloud-spread-tweets/cloud_run/recive/api"
)

var (
	PORT string
)

func init() {
	// 環境別の処理
	if runtime.GOOS == "linux" {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		PORT = fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
		log.Debug().Msgf("Linuxでの処理, PORT: %s", PORT)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		PORT = fmt.Sprintf("localhost:%s", os.Getenv("PORT"))
		log.Debug().Msgf("その他のOSでの処理, PORT: %s", PORT)
	}
}

func main() {
	key := os.Getenv("GOTWI_API_KEY_SECRET")
	secret := os.Getenv("GOTWI_API_KEY_SECRET")
	if key == "" || secret == "" {
		log.Fatal().Msg("X KEY, X SECRET is not set, does not load .env")
	}
	log.Debug().Msgf("app start, this consumer key: %s\n", strUpsideDown(key))

	e := echo.New()

	api.Routers(e)

	log.Fatal().Err(e.Start(PORT))
}

// strUpsideDown文字伏せ
func strUpsideDown(s string) string {
	l := len([]rune(s))
	if l <= 2 {
		return s
	}

	suffN := int(math.Floor(float64(l)/float64(2))) + 1
	suff := []rune(s)[suffN:]

	down := make([]string, suffN)
	for i := 0; i < suffN; i++ {
		down[i] = "*"
	}

	return strings.ReplaceAll(s, string(suff), strings.Join(down, ""))
}
