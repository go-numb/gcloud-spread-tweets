package api

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/go-numb/gcloud-spread-tweets/models"
)

type Client struct {
	ProjectID  string
	GUIURL     string
	POSTAPIURL string

	Key         string
	Secret      string
	CallbackURL string

	CredentialFile string

	// SheetID @MasterFile
	SpreadID    string
	SheetUserID string
	SheetPostID string
	RangeKey    string

	Firestore *models.ClientForFirestore
}

func New() *Client {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting working directory")
	}

	projectID := os.Getenv("PROJECTID")
	filename := filepath.Join(root, os.Getenv("CREDENTIALS"))
	return &Client{
		ProjectID:  projectID,
		GUIURL:     os.Getenv("GUIURL"),
		POSTAPIURL: os.Getenv("POSTAPIURL"),

		Key:         os.Getenv("GOTWI_API_KEY"),
		Secret:      os.Getenv("GOTWI_API_KEY_SECRET"),
		CallbackURL: os.Getenv("CALLBACK"),

		CredentialFile: filename,

		SpreadID:    os.Getenv("SPREADID"),
		SheetUserID: os.Getenv("SHEETUSER"),
		SheetPostID: os.Getenv("SHEETPOST"),
		RangeKey:    os.Getenv("RANGEKEY"),

		Firestore: &models.ClientForFirestore{
			ProjectID:      projectID,
			CredentialFile: filename,
		},
	}
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Routers(e *echo.Echo) {
	client := New()

	e.Use(
		middleware.CORS(),
		// session.Middleware(sessions.NewCookieStore([]byte(key))),
	)

	apiRouters := e.Group("/api")

	// - GET /api/health: health check
	apiRouters.GET("/health", client.Health)

	// - GET /api/data/usage: get usage data, data.users, data.posts
	apiRouters.GET("/data/usage", client.Usage)

	// - GET /api/upload/: registor update spreadsheet
	apiRouters.GET("/spreadsheet/upload", client.Registor)

	// Twitter/X callback
	// for Twitter API AccessToken,Secret
	apiRouters.GET("/x/auth/request", client.Auth)
	apiRouters.GET("/x/auth/callback", client.Callback)

	// for Database when cloud schedules
	// Account, Posts登録後のルーティン処理
	apiRouters.GET("/x/process", client.ScheduleProcess)

	// data control
	apiRouters.GET("/x/accounts", client.GetAccounts)

	// data control: account
	apiRouters.PUT("/x/account", client.PutAccount)

	// Google cloud schedulesからのトリガーで実行される
	// Post用データ取得・整形・実行Request
	apiRouters.GET("/x/posts", client.GetPosts)
	apiRouters.POST("/x/posts", client.CreatePosts)

	// Postデータ操作
	apiRouters.GET("/x/post", client.GetPost)       // 1件取得
	apiRouters.POST("/x/post", client.CreatePost)   // 新規作成
	apiRouters.PUT("/x/post", client.PutPost)       // 修正含む更新
	apiRouters.DELETE("/x/post", client.DeletePost) // 削除
}

func (p *Client) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "OK", Data: nil})
}

func (p *Client) Usage(c echo.Context) error {
	ctx := context.Background()
	store, err := p.Firestore.NewClient(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer store.Close()

	// get accounts
	docs, err := store.Collection(Accounts).
		Documents(ctx).GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	accounts := len(docs)

	docs, err = store.Collection(Posts).
		Documents(ctx).GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// get posts
	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "OK", Data: echo.Map{
		"users": accounts,
		"posts": len(docs),
	}})
}
