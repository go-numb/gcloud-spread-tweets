package models

import (
	"context"
	"log"

	sp "github.com/go-numb/go-spread-utils"

	"google.golang.org/api/sheets/v4"
)

type Client struct {
	Sheets *sp.Client
}

type Account struct {
}

type Post struct {
	AccessToken       string
	AccessTokenSecret string
}

func NewClient(ctx context.Context) *Client {
	serv, err := sheets.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		Sheets: &sp.Client{
			Sheets: serv,
		},
	}
}
