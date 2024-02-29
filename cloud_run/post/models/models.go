package models

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/go-numb/gcloud-spread-tweets/cloud_run/recive/libs"
	"golang.org/x/exp/rand"
)

type Client struct {
	Firestore *firestore.Client
}

type Account struct {
	libs.Account
}

type Post struct {
	AccessToken  string
	AccessSecret string

	libs.Post
}

func NewClient(ctx context.Context) *Client {
	config := &firebase.Config{ProjectID: os.Getenv("PROJECTID")}
	fire, err := firebase.NewApp(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	api, err := fire.Firestore(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		Firestore: api,
	}
}

// one 配列からランダムで一つ選択する
func one(l int) int {
	return rand.Intn(l)
}
