package libs

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/rs/zerolog/log"
)

func (p *Post) UploadMedias(filenames ...string) ([]string, error) {
	api := anaconda.NewTwitterApiWithCredentials(
		p.AccessToken,
		p.AccessSecret,
		os.Getenv("CUNSUMERKEY"),
		os.Getenv("CUNSUMERSECRET"),
	)

	var mediaIds []string
	for _, filename := range filenames {
		// ファイル拡張子で画像か動画か判定
		ext := strings.ToLower(filename[strings.LastIndex(filename, ".")+1:])
		if ext == "jpg" || ext == "jpeg" || ext == "png" || ext == "gif" || ext == "webp" {
			mediaId, err := UploadImage(api, filename)
			if err != nil {
				log.Error().Str("function", "UploadMedias").Err(err).Msgf("failed to upload image, %s", filename)
				continue
			}
			mediaIds = append(mediaIds, mediaId)
		} else if ext == "mp4" || ext == "mov" {
			mediaId, err := UploadVideo(api, filename)
			if err != nil {
				log.Error().Str("function", "UploadMedias").Err(err).Msgf("failed to upload video, %s", filename)
				continue
			}
			mediaIds = append(mediaIds, mediaId)
		}
	}

	if len(mediaIds) == 0 {
		return nil, fmt.Errorf("no media uploaded")
	}

	return mediaIds, nil
}

func UploadImage(api *anaconda.TwitterApi, filename string) (string, error) {
	base64Str, err := toBase64(filename)
	if err != nil {
		return "", err
	}
	media, err := api.UploadMedia(base64Str)
	if err != nil {
		return "", err
	}

	return media.MediaIDString, nil
}

func UploadVideo(api *anaconda.TwitterApi, filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		return "", err
	}

	fileType := "video/mp4"
	if strings.HasSuffix(filename, ".mov") {
		fileType = "video/quicktime"
	}

	media, err := api.UploadVideoInit(int(finfo.Size()), fileType)
	if err != nil {
		return "", err
	}

	log.Debug().Str("function", "UploadVideo").Msgf("video upload to twitter, %v, %d", fileType, finfo.Size())

	// チャンクサイズ（例：1MB）
	chunkSize := 1024 * 1024
	buffer := make([]byte, chunkSize)
	var segmentIndex int

	for {
		bytesRead, err := f.Read(buffer)
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("failed to read file, %v", err)
		}
		if bytesRead == 0 {
			break
		}

		if err := api.UploadVideoAppend(media.MediaIDString, segmentIndex, base64.StdEncoding.EncodeToString(buffer[:bytesRead])); err != nil {
			return "", fmt.Errorf("failed to upload video append, %v", err)
		}
		segmentIndex++
	}

	// アップロードの完了
	result, err := api.UploadVideoFinalize(media.MediaIDString)
	if err != nil {
		return "", fmt.Errorf("failed to upload video finalize, %v", err)
	}

	time.Sleep(10 * time.Second)
	log.Debug().Str("function", "UploadVideo").Msgf("video uploaded to twitter, %s, %v, %d", media.MediaIDString, result.Video.VideoType, result.Size)

	return media.MediaIDString, nil
}

func toBase64(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
