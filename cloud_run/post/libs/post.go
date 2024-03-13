package libs

import (
	"context"
	"fmt"

	"github.com/go-numb/gcloud-spread-tweets/models"
	xpostblue "github.com/go-numb/x-post-to-blue"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
)

const XLIMIT = 141

type Post struct {
	AccessToken  string
	AccessSecret string
	Password     string

	*models.Post
}

func (p *Client) ToDecrypto(account models.Account, post *Post) error {
	decPassword, err := models.DecryptPassword(account.Password, p.PasswordController)
	if err != nil {
		return err
	}
	post.AccessToken = account.AccessToken
	post.AccessSecret = account.AccessSecret
	post.Password = decPassword

	return nil
}

func (p *Post) Do() error {
	if !IsBlue(p.Text) {
		return p.DoAPI()
	}

	return p.DoGUI()
}

func (p *Post) DoAPI() error {
	op := &types.CreateInput{}
	medias, isMedia := isMedia(p)
	if isMedia {
		mediaIDs, err := UploadMedias(p.AccessToken, p.AccessSecret, medias...)
		if err != nil {
			// 画像Uploadでエラーが出たときは、WithFiles設定に従う
			if p.WithFiles == 1 {
				return err
			}

			// ファイルがなくもテキストのみで投稿へ移行する
		} else {
			op.Media = &types.CreateInputMedia{
				MediaIDs: mediaIDs,
			}
		}
	}

	client, err := new(p.AccessToken, p.AccessSecret)
	if err != nil {
		return err
	}
	res, err := managetweet.Create(context.Background(), client, op)
	if err != nil {
		return err
	} else if res.HasPartialError() {
		return fmt.Errorf("create post partial error")
	}

	return nil
}

// TODO: GUI for blue
func (p *Post) DoGUI() error {
	isHeadless := true
	client := xpostblue.New(isHeadless)
	defer client.Close()

	if err := client.Login(p.ID, p.Password); err != nil {
		return err
	}

	medias, _ := isMedia(p)
	if err := client.Post(true, 60, p.Text, medias...); err != nil {
		return err
	}

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

func isMedia(post *Post) ([]string, bool) {
	var results []string
	medias := []string{post.File1, post.File2, post.File3, post.File4}
	for _, media := range medias {
		if media != "" {
			results = append(results, media)
		}
	}

	if len(results) == 0 {
		return results, false
	}

	return results, true
}

// 文字数での判定
func IsBlue(s string) bool {
	if len([]rune(s)) < XLIMIT {
		return false
	}
	return true
}
