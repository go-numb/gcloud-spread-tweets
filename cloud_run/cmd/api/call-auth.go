package api

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"github.com/labstack/echo/v4"
)

func (p *Client) Auth(c echo.Context) error {
	log.Debug().Msgf("call Auth")

	var config = oauth1.Config{
		ConsumerKey:    p.Key,                     // アプリのConsumer Key
		ConsumerSecret: p.Secret,                  // アプリのConsumer Secret
		CallbackURL:    p.CallbackURL,             // コールバックURL
		Endpoint:       twitter.AuthorizeEndpoint, // TwitterのOAuthエンドポイント
	}

	// リクエストトークンと認証URLの取得
	requestToken, requestTokenSecret, err := config.RequestToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error getting request token, %v", err)})
	}

	claims := Claims{
		RequestToken:       requestToken,
		RequestTokenSecret: requestTokenSecret,
	}
	if err := p.SetFirestore(claims); err != nil {
		log.Debug().Msgf("error setting firestore, %v", err)
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting firestore, %v", err)})
	}

	// Twitterの認証ページへリダイレクト
	authURL, err := config.AuthorizationURL(requestToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: "Error getting authorization URL"})
	}

	value, err := p.GetFirestore(claims.RequestToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error getting firestore, %v", err)})
	}

	log.Debug().Msgf("token: %+v", value)

	return c.JSON(http.StatusOK, Response{Code: http.StatusOK, Message: "Success", Data: map[string]string{"url": authURL.String(), "token": requestToken}})
}

// Callback is a handler for the OAuth callback
func (p *Client) Callback(c echo.Context) error {
	log.Debug().Msgf("call Callback")

	// リクエストトークンを取得
	requestToken := c.QueryParam("oauth_token")
	verifier := c.QueryParam("oauth_verifier")

	claims, err := p.GetFirestore(requestToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error getting firestore, %v", err)})
	}

	log.Debug().Msgf("request token: %s, secret: %s, verifier: %s", requestToken, claims.RequestTokenSecret, verifier)

	var config = oauth1.Config{
		ConsumerKey:    p.Key,                     // アプリのConsumer Key
		ConsumerSecret: p.Secret,                  // アプリのConsumer Secret
		CallbackURL:    p.CallbackURL,             // コールバックURL
		Endpoint:       twitter.AuthorizeEndpoint, // TwitterのOAuthエンドポイント
	}

	accessToken, accessSecret, err := config.AccessToken(requestToken, claims.RequestTokenSecret, verifier)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error getting access token, %v", err)})
	}

	claims.AccessToken = accessToken
	claims.AccessSecret = accessSecret
	if err := p.SetFirestore(*claims); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("Error setting firestore, %v", err)})
	}

	log.Debug().Msgf("end callback access token: %s, secret: %s\n", accessToken, accessSecret)

	// ページ変異先にリダイレクト
	return c.Redirect(http.StatusTemporaryRedirect, p.GUIURL)
}
