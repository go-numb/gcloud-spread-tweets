package api

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/go-numb/gcloud-spread-tweets/models"
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
		ProjectID: projectID,
		GUIURL:    os.Getenv("GUIURL"),

		Key:         os.Getenv("KEY"),
		Secret:      os.Getenv("SECRET"),
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
	// - GET /api/upload/: registor update spreadsheet
	apiRouters.GET("/spreadsheet/upload", client.Registor)
	// Twitter/X callback
	// for Twitter API AccessToken,Secret
	apiRouters.GET("/x/auth/request", client.Auth)
	apiRouters.GET("/x/auth/callback", client.Callback)

	// for Database
	// Account, Posts登録後のルーティン処理
	apiRouters.GET("/x/accounts", GetAccounts)
	// Google cloud schedulesからのトリガーで実行される
	// Post用データ取得・整形・実行Request
	apiRouters.GET("/x/post", GetPost)

	// Postデータ操作
	apiRouters.GET("/x/post", GetPost)       // 1件取得
	apiRouters.POST("/x/post", CreatePost)   // 新規作成
	apiRouters.PUT("/x/post", PutPost)       // 修正含む更新
	apiRouters.DELETE("/x/post", DeletePost) // 削除
}

func GetAccounts(c echo.Context) error {
	return c.JSON(200, "accounts")
}

// GetPost
func GetPosts(c echo.Context) error {
	return c.JSON(200, "pong")
}

func GetPost(c echo.Context) error {
	return c.JSON(200, "pong")
}

func CreatePost(c echo.Context) error {
	return c.JSON(200, "pong")
}

func PutPost(c echo.Context) error {
	return c.JSON(200, "pong")
}

func DeletePost(c echo.Context) error {
	return c.JSON(200, "pong")
}
