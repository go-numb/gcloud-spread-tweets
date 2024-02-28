package api

import (
	"context"
	"fmt"

	"go-api/libs"
	"net/http"
	"strings"

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
	claims := libs.Claims{}
	if err := p.GetAnyFirestore(Users, customerToken, &claims); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error getting session or nothing token: %s in session, %v", customerToken, err)})
	}

	// アクセストークンの確認
	if claims.AccessToken == "" || claims.AccessSecret == "" {
		return c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: "Error getting access token"})
	}

	log.Debug().Msgf("customers post SpreadsheetID: %s, access token: %s, secret token: %s", customerSpreadsheetID, claims.AccessToken, claims.AccessSecret)

	// Spreadsheet API Clientの初期化
	client := spreads.New(
		context.Background(),
		p.CredentialFile,
	)

	// 顧客、登録Spreadsheet:postsのデータの読み込み
	client.
		SetSpreadID(customerSpreadsheetID).
		SetSheetName(p.SheetPostID).
		SetRangeKey(p.RangeKey)
	customerPosts := []libs.Post{}
	customerPostDf, err := client.Read(&customerPosts)
	if err != nil || len(customerPosts) == 0 {
		log.Debug().Err(err).Msg("Error reading customer tweets spreadsheet")
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: "Error reading customer tweets spreadsheet"})
	}

	// 顧客登録Spreadsheetのデータ型式を確認
	// 顧客データのカラム名とカラム数が一致しているかを確認
	if err := libs.CheckColumns(customerPostDf); err != nil {
		log.Debug().Err(err).Msg("Error checking columns")
		return c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error checking columns > %v", err)})
	}

	log.Printf("checked customers sheet, rows: %d, columns: %v", customerPostDf.Nrow(), customerPostDf.Names())

	// マスターファイルのアカウントデータの読み込み
	customerAccountIds := libs.GetUniqueKeys(customerPosts, func(t libs.Post) string {
		return t.ID
	})
	// 重複チェック
	customerAccountIds = libs.CheckDuplicate(customerAccountIds)
	if err := p.CheckExistKeysFirestore(Accounts, customerAccountIds); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error checking customer account keys > %v", err)})
	}

	// 顧客のPostsデータをマスターファイルへ追記
	if err := p.SetAnyFirestore(Posts, "", customerPosts); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting posts firestore > %v", err)})
	}

	// master: x_usersの末尾に追記
	// Twitter/Xアカウントで認証し、そのアカウントを含むx_postsのデータを登録する
	// そのため、認証アカウントとx_posts内のアカウントが一致している
	account := libs.NewAccount(
		customerPosts[0].ID,
		customerSpreadsheetID,
		claims.AccessToken,
		claims.AccessSecret).
		// Default: Free Plan
		SetSubscribed(libs.SubscribedFree).
		SetTime([]int{17}, []int{0}).
		SetTerm(48)
	if err := p.SetAnyFirestore(Accounts, account.ID, account); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting account firestore > %v", err)})
	}

	// SessionにXUIDをセット
	claims.Name = customerPosts[0].ID
	claims.SpreadsheetID = customerSpreadsheetID
	if err := p.SetAnyFirestore(Users, claims.RequestToken, claims); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting firestore > %v", err)})
	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: claims})
}
