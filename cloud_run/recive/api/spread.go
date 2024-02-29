package api

import (
	"context"
	"fmt"

	"net/http"
	"strings"

	models "github.com/go-numb/gcloud-spread-tweets/cloud_run/models"

	"github.com/go-gota/gota/dataframe"
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
func (p *Client) Regi(c echo.Context) error {
	log.Debug().Msgf("call Registor")

	// フォームからのPOSTを表記する
	customerSpreadsheetID := strings.TrimSpace(c.QueryParam("spreadsheet_id"))
	customerToken := strings.TrimSpace(c.QueryParam("token"))
	if customerSpreadsheetID == "" {
		return c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: "Error getting query for customer spreadsheet id"})
	} else if customerToken == "" {
		return c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: "Error getting query for customer token"})
	}

	claims := models.Claims{}
	if err := p.GetAnyFirestore(Users, customerToken, &claims); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error getting session, %v", err)})
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
	customerPosts := []models.Post{}
	customerPostDf, err := client.Read(&customerPosts)
	if err != nil || len(customerPosts) == 0 {
		log.Error().Err(err).Msg("Error reading customer tweets spreadsheet")
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: "Error reading customer tweets spreadsheet"})
	}

	// 顧客登録Spreadsheetのデータ型式を確認
	// 顧客データのカラム名とカラム数が一致しているかを確認
	if err := models.CheckColumns(customerPostDf); err != nil {
		log.Error().Err(err).Msg("Error checking columns")
		return c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error checking columns > %v", err)})
	}

	log.Printf("checked customers sheet, rows: %d, columns: %v\n", customerPostDf.Nrow(), customerPostDf.Names())

	// マスターファイルのユーザーデータの読み込み
	client.SetSpreadID(p.SpreadID).SetSheetName(p.SheetUserID)
	masterCustomerAccounts := []models.Account{}
	masterCustomerAccountsDf, err := client.Read(&masterCustomerAccounts)
	if err != nil {
		log.Error().Err(err).Msg("Error reading master users spreadsheet")
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: "Error reading master users spreadsheet"})
	}

	// 登録SpreadsheetID及びTwitter/Xアカウントの重複チェックがマスターファイルにすでに存在しないかチェック
	if err := models.CheckDupID(customerPosts[0].ID, customerSpreadsheetID, masterCustomerAccounts); err != nil {
		log.Error().Err(err).Msg("Error checking duplicate id")
		return c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error checking duplicate id > %v", err)})
	}

	//
	// master x_postsのデータの読み込み
	client.SetSpreadID(p.SpreadID).SetSheetName(p.SheetPostID)
	masterPosts := []models.Post{}
	masterPostDf, err := client.Read(&masterPosts)
	if err != nil || len(masterPosts) == 0 {
		log.Error().Err(err).Msg("Error reading master tweets spreadsheet")
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: "Error reading customer tweets spreadsheet"})
	}

	// 顧客のPostsデータをマスターファイルposts末尾へ追記し更新
	customerRecords := customerPostDf.Records()
	// 1行目はカラム名のため、2行目からデータを取得
	updateRecords := spreads.ConvertStringToInterface(customerRecords)[1:]
	// master: x_postsの末尾に追記
	rowN := models.DfNrowToLastNrow(masterPostDf)
	client.SetRangeKey(fmt.Sprintf("A%d:Z", rowN)) // 末尾に追記
	if err := client.Update(updateRecords); err != nil {
		log.Error().Err(err).Msg("Error updating master tweets spreadsheet")
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error updating master tweets spreadsheet > %v", err)})
	}

	// master: x_usersの末尾に追記
	// Twitter/Xアカウントで認証し、そのアカウントを含むx_postsのデータを登録する
	// そのため、認証アカウントとx_posts内のアカウントが一致している
	new_account := []models.Account{
		*models.NewAccount(
			customerPosts[0].ID,
			customerSpreadsheetID,
			claims.AccessToken,
			claims.AccessSecret).
			// Default: Free Plan
			SetSubscribed(models.SubscribedFree).
			SetTime([]int{17}, []int{0}).
			SetTerm(48),
	}
	temp := dataframe.LoadStructs(new_account)
	records := temp.Subset(0).Records()
	updateRow := spreads.ConvertStringToInterface(records)[1]
	client.SetSheetName(p.SheetUserID) // master: x_usersに追記
	rowN = models.DfNrowToLastNrow(masterCustomerAccountsDf)
	client.SetRangeKey(fmt.Sprintf("A%d:Z%d", rowN, rowN)) // 末尾に追記
	if err := client.UpdateRow(updateRow); err != nil {
		log.Error().Err(err).Msg("Error updating master users spreadsheet")
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error updating master users spreadsheet > %v", err)})
	}

	// SessionにXUIDをセット
	claims.Name = customerPosts[0].ID
	claims.SpreadsheetID = customerSpreadsheetID

	if err := p.SetAnyFirestore(Users, claims.RequestToken, claims); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting firestore > %v", err)})
	}

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: claims})
}
