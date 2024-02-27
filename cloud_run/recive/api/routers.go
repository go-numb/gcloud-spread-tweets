package api

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Client struct {
	ProjectID string
	GUIURL    string

	Key         string
	Secret      string
	CallbackURL string

	CredentialFile string

	// SheetID @MasterFile
	SpreadID    string
	SheetUserID string
	SheetPostID string
	RangeKey    string
}

func New() *Client {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting working directory")
	}

	return &Client{
		ProjectID: os.Getenv("PROJECTID"),
		GUIURL:    os.Getenv("GUIURL"),

		Key:         os.Getenv("KEY"),
		Secret:      os.Getenv("SECRET"),
		CallbackURL: os.Getenv("CALLBACK"),

		CredentialFile: filepath.Join(root, os.Getenv("CREDENTIALS")),

		SpreadID:    os.Getenv("SPREADID"),
		SheetUserID: os.Getenv("SHEETUSER"),
		SheetPostID: os.Getenv("SHEETPOST"),
		RangeKey:    os.Getenv("RANGEKEY"),
	}
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Routers(e *echo.Echo) {
	client := New()

	key := os.Getenv("SESSIONKEY")
	log.Debug().Msgf("session key: %s", key)
	e.Use(
		middleware.CORS(),
		// session.Middleware(sessions.NewCookieStore([]byte(key))),
	)

	apiRoutes := e.Group("/api")
	// - GET /api/upload/: registor update spreadsheet
	apiRoutes.GET("/spreadsheet/upload", client.Regi)
	// Twitter/X callback
	// for Twitter API AccessToken,Secret
	apiRoutes.GET("/x/auth/request", client.Auth)
	apiRoutes.GET("/x/auth/callback", client.Callback)
}
