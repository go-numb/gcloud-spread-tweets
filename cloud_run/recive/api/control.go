package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/go-numb/gcloud-spread-tweets/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// GetAccounts 現在時刻に合致したアカウントを取得
func (p *Client) GetAccounts(c echo.Context) error {
	ctx := context.Background()
	accounts, err := p.getAccounts(ctx, time.Now())
	if err != nil {
		log.Debug().Err(err).Msg("failed to get accounts")
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: accounts})
}

// PutAccount アカウントを更新
func (p *Client) PutAccount(c echo.Context) error {
	token := c.QueryParam("token")
	username := c.QueryParam("username")

	ctx := context.Background()
	if err := p.IsUserChecker(ctx, token, username); err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	var account models.Account
	if err := c.Bind(&account); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	fmt.Printf("%+v\n", account)

	store, err := p.Firestore.NewClient(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting firestore, %v", err)})
	}
	defer store.Close()

	encPassword, err := models.EncryptPassword(account.Password, p.PasswordController)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error encrypting password, %v", err)})
	}
	account.Password = ""

	if _, err := store.Collection(Accounts).Doc(username).Update(ctx, []firestore.Update{
		{Path: "password", Value: encPassword},
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error updating firestore, %v", err)})
	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: account})
}

// GetPost Postを取得 query: token, username
func (p *Client) GetPosts(c echo.Context) error {
	token := c.QueryParam("token")
	username := c.QueryParam("username")

	ctx := context.Background()
	if err := p.IsUserChecker(ctx, token, username); err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	posts, err := p.getPosts(ctx, models.Account{
		ID: username,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: posts})
}

// CreatePosts Postsを新規作成 query: spread_id
func (p *Client) CreatePosts(c echo.Context) error {
	token := c.QueryParam("token")
	username := c.QueryParam("username")

	ctx := context.Background()
	if err := p.IsUserChecker(ctx, token, username); err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	spreadID := c.QueryParam("spread_id")
	if spreadID == "" {
		return c.JSON(http.StatusBadRequest, "spread_id is empty")
	}

	// ここでspreadsheetからデータを取得
	// カラムのチェックを行い、データを取得
	posts, err := p.ReadSpreadsheetForPost(spreadID)
	if err != nil {
		log.Error().Err(err).Msg("Error reading customer tweets spreadsheet")
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error reading customer tweets spreadsheet > %v", err)})
	}

	if err := p.Firestore.Set(ctx, Posts, "", posts); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting posts firestore > %v", err)})
	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: posts})
}

// GetPost Postを取得 query: post_id
func (p *Client) GetPost(c echo.Context) error {
	token := c.QueryParam("token")
	username := c.QueryParam("username")

	ctx := context.Background()
	if err := p.IsUserChecker(ctx, token, username); err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	id := c.QueryParam("post_id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: "post_id is empty"})
	}

	var post models.Post
	if err := p.Firestore.Get(ctx, Posts, id, &post); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error getting firestore, %v", err)})
	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: post})
}

// CreatePost Postを新規作成 from form values
func (p *Client) CreatePost(c echo.Context) error {
	token := c.QueryParam("token")
	username := c.QueryParam("username")

	ctx := context.Background()
	if err := p.IsUserChecker(ctx, token, username); err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	// フォームデータからデータを取得
	// １つの投稿データを保存
	var post models.Post
	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: fmt.Sprintf("invalid request, %v", err)})
	}

	// Check
	if post.ID == "" {
		return c.JSON(http.StatusBadRequest, "id is empty")
	} else if post.Text == "" {
		return c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: "text is empty"})
	} else if len([]byte(post.Text)) > 280 {
		return c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: "text is over 280 characters, please shorten it or subscribe to the premium plan."})
	}

	post.UUID = uuid.New().String()
	post.SetCreateAt()

	if err := p.Firestore.Set(ctx, Posts, post.UUID, post); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting firestore, %v", err)})
	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: post})
}

// PutPost Postを更新
// 主にchecked, priorityの更新(削除除く)
func (p *Client) PutPost(c echo.Context) error {
	token := c.QueryParam("token")
	username := c.QueryParam("username")

	ctx := context.Background()
	if err := p.IsUserChecker(ctx, token, username); err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	var post models.Post
	if err := c.Bind(&post); err != nil {
		// 失敗したときにリクエストボディをStringとして出力する
		log.Error().Err(err).Msg("failed to bind request body")
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	app, err := p.Firestore.NewClient(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting firestore, %v", err)})
	}
	defer app.Close()

	// 指定項目のみ更新を許可する
	// Why: 他の項目が更新されると、他の項目が空になる可能性や、不正な値が入る可能性があるため
	// GUI: 値はフォームに入力済み、変更し血の更新をリクエストする。よって、空になっていれば、空になることを許可する
	updateValues := []firestore.Update{
		{Path: "file_1", Value: post.File1},
		{Path: "file_2", Value: post.File2},
		{Path: "file_3", Value: post.File3},
		{Path: "file_4", Value: post.File4},
		{Path: "with_files", Value: post.WithFiles},
		{Path: "checked", Value: post.Checked},
		{Path: "priority", Value: post.Priority},
	}
	if post.Text != "" {
		updateValues = append(updateValues, firestore.Update{Path: "text", Value: post.Text})
	}
	if _, err := app.Collection(Posts).Doc(post.UUID).Update(ctx, updateValues); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error updating firestore, %v", err)})

	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: post})
}

// DeletePost Postを削除
func (p *Client) DeletePost(c echo.Context) error {
	token := c.QueryParam("token")
	username := c.QueryParam("username")

	ctx := context.Background()
	if err := p.IsUserChecker(ctx, token, username); err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	uuid := c.QueryParam("uuid")
	if uuid == "" {
		return c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: "uuid is empty"})
	}

	app, err := p.Firestore.NewClient(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting firestore, %v", err)})
	}
	defer app.Close()

	if _, err := app.Collection(Posts).Doc(uuid).Update(ctx, []firestore.Update{
		{Path: "is_delete", Value: true},
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error deleting firestore, %v", err)})
	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: uuid})
}

func (p *Client) IsUserChecker(ctx context.Context, token, username string) error {
	if token == "" || username == "" {
		return fmt.Errorf("token or username is empty")
	}

	if !p.isOK(ctx, token, username) {
		return fmt.Errorf("token or username is invalid")
	}

	return nil
}
