package main

import (
	"context"
	"os"

	aiozstreamsdk "github.com/AIOZNetwork/aioz-stream-go-client"
	"github.com/joho/godotenv"
)

func countVideoAndAudio(ctx context.Context) (int, int, error) {
	errEnv := godotenv.Load()
	if errEnv != nil {
		return 0, 0, errEnv
	}

	publicKey := os.Getenv("PUBLIC_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	apiCreds := aiozstreamsdk.AuthCredentials{
		PublicKey: publicKey,
		SecretKey: secretKey,
	}

	client := aiozstreamsdk.ClientBuilder(apiCreds).Build()

	videoReq := aiozstreamsdk.NewGetMediaListRequest()
	videoReq.SetType("video")

	videoRes, err := client.Media.GetMediaListWithContext(ctx, *videoReq)
	if err != nil {
		return 0, 0, err
	}

	videoCount := int(*videoRes.GetData().Total)

	audioReq := aiozstreamsdk.NewGetMediaListRequest()
	audioReq.SetType("audio")

	audioRes, err := client.Media.GetMediaListWithContext(ctx, *audioReq)
	if err != nil {
		return 0, 0, err
	}

	audioCount := int(*audioRes.GetData().Total)

	return videoCount, audioCount, nil
}
