package api

import (
	"context"
	"fmt"

	"net/http"
	"strings"

	"github.com/go-numb/gcloud-spread-tweets/models"

	spreads "github.com/go-numb/go-spread-utils"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// Registor handles the form submission
/*
1. フロントエンドからSpreadsheetIDを取得
2. セッションからアクセストークンを取得
3. 顧客のSpreadsheetIDの投稿データを読み込み
4. 顧客のSpreadsheetIDの投稿データ型式を確認
5. マスターファイルのユーザーデータを読み込み
6. 顧客SpreadsheetID及びTwitter/Xアカウントの重複チェックがマスターファイルにすでに存在しないかチェック
7. マスターファイルの投稿データを取得
8. マスターファイルの投稿データ末尾へ追記し更新
*/
func (p *Client) Registor(c echo.Context) error {
	log.Debug().Msgf("call Registor")

	// フォームからのPOST[spread_id, token]を受け取る
	customerSpreadsheetID := strings.TrimSpace(c.QueryParam("spreadsheet_id"))
	customerToken := strings.TrimSpace(c.QueryParam("token"))
	if customerSpreadsheetID == "" {
		return c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: "Error getting query for customer spreadsheet id"})
	} else if customerToken == "" {
		return c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: "Error getting query for customer token"})
	}

	// アクセストークンをキーにしてセッションを取得
	// セッションが有効か確認
	ctx := context.Background()
	claims := models.Claims{}
	if err := p.Firestore.Get(ctx, Users, customerToken, &claims); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error getting session or nothing token: %s in session, %v", customerToken, err)})
	}

	// アクセストークンの確認
	if claims.AccessToken == "" || claims.AccessSecret == "" {
		return c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: "Error getting access token"})
	}

	log.Debug().Msgf("customers post SpreadsheetID: %s, access token: %s, secret token: %s", customerSpreadsheetID, claims.AccessToken, claims.AccessSecret)

	// Spreadsheet API Clientの初期化
	posts, err := p.ReadSpreadsheetForPost(customerSpreadsheetID)
	if err != nil {
		log.Error().Err(err).Msg("Error reading customer tweets spreadsheet")
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error reading customer tweets spreadsheet > %v", err)})
	}

	// Auth認証を行ったAccountIDと一致する投稿データを取得
	var customerPosts []models.Post
	for _, post := range posts {
		if claims.ID != post.ID {
			continue
		}
		customerPosts = append(customerPosts, post)
	}

	if _, err := p.Firestore.IsExist(ctx, Accounts, claims.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error checking account exist > %v", err)})
	}

	// 顧客のPostsデータをマスターファイルへ追記
	if err := p.Firestore.Set(ctx, Posts, "", customerPosts); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting posts firestore > %v", err)})
	}

	// master: x_usersの末尾に追記
	// Twitter/Xアカウントで認証し、そのアカウントを含むx_postsのデータを登録する
	// そのため、認証アカウントとx_posts内のアカウントが一致している
	account := models.NewAccount(
		customerPosts[0].ID,
		customerSpreadsheetID,
		claims.AccessToken,
		claims.AccessSecret).
		// Default: Free Plan
		SetSubscribed(models.SubscribedFree).
		SetTime([]int{17}, []int{0}).
		SetTerm(48)
	if err := p.Firestore.Set(ctx, Accounts, account.ID, account); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting account firestore > %v", err)})
	}

	// SessionにXUIDをセット
	claims.ID = customerPosts[0].ID
	claims.SpreadsheetID = customerSpreadsheetID
	if err := p.Firestore.Set(ctx, Users, claims.RequestToken, claims); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting firestore > %v", err)})
	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: claims})
}

func (p *Client) ReadSpreadsheetForPost(id string) ([]models.Post, error) {
	// Spreadsheet API Clientの初期化
	client := spreads.New(
		context.Background(),
		p.CredentialFile,
	)

	// 顧客、登録Spreadsheet:postsのデータの読み込み
	client.
		SetSpreadID(id).
		SetSheetName(p.SheetPostID).
		SetRangeKey(p.RangeKey)
	customerPosts := []models.Post{}
	customerPostDf, err := client.Read(&customerPosts)
	if err != nil || len(customerPosts) == 0 {
		return nil, fmt.Errorf("error reading customer tweets spreadsheet, length: %d", len(customerPosts))
	}

	// 顧客登録Spreadsheetのデータ型式を確認
	// 顧客データのカラム名とカラム数が一致しているかを確認
	if err := models.CheckColumns(customerPostDf); err != nil {
		return nil, fmt.Errorf("error checking columns > %v", err)
	}

	log.Printf("checked customers sheet, rows: %d, columns: %v", customerPostDf.Nrow(), customerPostDf.Names())

	return customerPosts, nil
}
