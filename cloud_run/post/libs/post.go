package libs

import (
	"context"
	"fmt"

	xpostblue "github.com/go-numb/x-post-to-blue"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
)

func (p *Post) Do() error {
	client, err := new(p.AccessToken, p.AccessSecret)
	if err != nil {
		return err
	}
	op := &types.CreateInput{}

	medias, isMedia := isMedia(p)
	if !isMedia {
		mediaIDs, err := p.UploadMedias(medias...)
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
