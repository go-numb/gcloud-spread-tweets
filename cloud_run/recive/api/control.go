package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-numb/gcloud-spread-tweets/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// GetPost Postを取得 query: token, username
func (p *Client) GetPosts(c echo.Context) error {
	token := c.QueryParam("token")
	username := c.QueryParam("username")
	if token == "" || username == "" {
		return c.JSON(http.StatusForbidden, "token or username is empty")
	}

	ctx := context.Background()
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

	ctx := context.Background()
	if err := p.Firestore.Set(ctx, Posts, "", posts); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting posts firestore > %v", err)})
	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: posts})
}

// GetPost Postを取得 query: post_id
func (p *Client) GetPost(c echo.Context) error {
	id := c.QueryParam("post_id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "id is empty")
	}

	ctx := context.Background()
	var post models.Post
	if err := p.Firestore.Get(ctx, Posts, id, &post); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error getting firestore, %v", err)})
	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: post})
}

// CreatePost Postを新規作成 from form values
func (p *Client) CreatePost(c echo.Context) error {
	// フォームデータからデータを取得
	// １つの投稿データを保存
	var post models.Post
	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	// Check
	if post.ID == "" {
		return c.JSON(http.StatusBadRequest, "id is empty")
	} else if post.Text == "" {
		return c.JSON(http.StatusBadRequest, "text is empty")
	} else if len([]byte(post.Text)) > 280 {
		return c.JSON(http.StatusBadRequest, "text is over 280 characters, please shorten it or subscribe to the premium plan.")
	}

	post.UUID = uuid.New().String()
	post.SetCreateAt()

	ctx := context.Background()
	if err := p.Firestore.Set(ctx, Posts, post.UUID, post); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting firestore, %v", err)})
	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: post})
}

// PutPost Postを更新
// 主にchecked, priorityの更新(削除除く)
func (p *Client) PutPost(c echo.Context) error {
	var post models.Post
	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	ctx := context.Background()
	if err := p.Firestore.Set(ctx, Posts, post.UUID, post); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting firestore, %v", err)})
	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: post})
}

// DeletePost Postを削除
func (p *Client) DeletePost(c echo.Context) error {
	var post models.Post
	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	post.IsDelete = true

	ctx := context.Background()
	if err := p.Firestore.Set(ctx, Posts, post.UUID, post); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting firestore, %v", err)})
	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: post})
}
