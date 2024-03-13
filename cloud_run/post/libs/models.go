package libs

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/go-numb/gcloud-spread-tweets/models"
	"golang.org/x/exp/rand"

	_ "github.com/joho/godotenv/autoload"
)

type Client struct {
	Firestore *models.ClientForFirestore

	PasswordController []byte
}

func NewClient() *Client {
	// キーを環境変数から取得 (base64エンコードされている)
	keystr := os.Getenv("PASSWORDCONTROLLER")
	// キーをbase64デコード
	key, err := base64.StdEncoding.DecodeString(keystr)
	if err != nil {
		log.Fatalf("Error decoding key, %v", err)
	}

	return &Client{
		Firestore: &models.ClientForFirestore{
			ProjectID: os.Getenv("PROJECTID"),
		},
		PasswordController: key,
	}
}

// one 配列からランダムで一つ選択する
func one(l int) int {
	return rand.Intn(l)
}
