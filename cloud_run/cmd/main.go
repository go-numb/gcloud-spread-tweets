package main

import (
	"math"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/joho/godotenv/autoload"

	"go-api/api"
)

const (
	IS_PRODUCTION = true
)

func init() {
	if IS_PRODUCTION {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

}

func main() {
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	if key == "" || secret == "" {
		log.Fatal().Msg("X KEY, X SECRET is not set, does not load .env")
	}
	log.Debug().Msgf("app start, this consumer key: %s\n", strUpsideDown(key))

	e := echo.New()

	api.Routers(e)

	log.Fatal().Err(e.Start(":" + os.Getenv("PORT")))
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
