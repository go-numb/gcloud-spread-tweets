package libs

import (
	"os"

	"github.com/go-numb/gcloud-spread-tweets/models"
	"golang.org/x/exp/rand"

	_ "github.com/joho/godotenv/autoload"
)

type Client struct {
	Firestore *models.ClientForFirestore
}

func NewClient() *Client {
	return &Client{
		Firestore: &models.ClientForFirestore{
			ProjectID: os.Getenv("PROJECTID"),
		},
	}
}

// one 配列からランダムで一つ選択する
func one(l int) int {
	return rand.Intn(l)
}
