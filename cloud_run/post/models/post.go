package models

import (
	"context"
	"fmt"

	xpostblue "github.com/go-numb/x-post-to-blue"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"

	"github.com/go-gota/gota/dataframe"
)

func (p *Client) GetAccounts(binder any) (dataframe.DataFrame, error) {
	p.Sheets.SetSpreadID("1").SetSheetName("users").SetRangeKey("A1")
	return p.Sheets.Read(binder)
}

func (p *Client) GetPosts(account Account) ([]Post, error) {
	return []Post{}, nil
}

func (p *Post) Do() error {
	client, err := new(p.AccessToken, p.AccessTokenSecret)
	if err != nil {
		return err
	}

	res, err := managetweet.Create(context.Background(), client, &types.CreateInput{})
	if err != nil {
		return err
	} else if res.HasPartialError() {
		return fmt.Errorf("create post partial error")
	}

	return nil
}

func (p *Post) DoGUI() error {
	isHeadless := true
	client := xpostblue.New(isHeadless)
	defer client.Close()

	return nil
}

func new(accesstoken, accesssecret string) (*gotwi.Client, error) {
	config := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           accesstoken,
		OAuthTokenSecret:     accesssecret,
	}

	return gotwi.NewClient(config)
}
