package tool

import (
	"context"
	"mcp-weather-server/model"

	aiozstreamsdk "github.com/AIOZNetwork/aioz-stream-go-client"
)

func createClient(publicKey string, secretKey string) *aiozstreamsdk.Client {
	apiCreds := aiozstreamsdk.AuthCredentials{
		PublicKey: publicKey,
		SecretKey: secretKey,
	}

	client := aiozstreamsdk.ClientBuilder(apiCreds).Build()
	return client
}

func CountVideoAndAudio(ctx context.Context, publicKey string, secretKey string) (int, int, error) {
	client := createClient(publicKey, secretKey)

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

func GetVideoDetailByName(ctx context.Context, publicKey string, secretKey string, videoName string) (interface{}, error) {
	client := createClient(publicKey, secretKey)
	videoReq := aiozstreamsdk.NewGetMediaListRequest()
	videoReq.SetType("video")
	videoReq.SetSearch(videoName)

	videoRes, err := client.Media.GetMediaListWithContext(ctx, *videoReq)
	if err != nil {
		return nil, err
	}

	sourceURLs := make(map[string]string)

	media := videoRes.GetData().Media
	if media == nil {
		return sourceURLs, nil
	}

	for _, m := range *media {
		sourceURLs["EmbededURL"] = *m.Assets.DashPlayerUrl
		sourceURLs["Mp4URL"] = *m.Assets.Mp4Url
		sourceURLs["Thumbnail"] = *m.Assets.ThumbnailUrl
		sourceURLs["SourceURL"] = *m.Assets.SourceUrl
	}

	return sourceURLs, nil
}

func GetVideos(ctx context.Context, publicKey string, secretKey string) (interface{}, error) {
	client := createClient(publicKey, secretKey)
	videoReq := aiozstreamsdk.NewGetMediaListRequest()
	videoReq.SetType("video")

	videoRes, err := client.Media.GetMediaListWithContext(ctx, *videoReq)
	if err != nil {
		return nil, err
	}

	var mediaList []model.MediaInfo
	media := videoRes.GetData().Media
	if media == nil {
		return mediaList, nil
	}

	for _, m := range *media {
		item := model.MediaInfo{}

		if m.Id != nil {
			item.MediaID = *m.Id
		}

		if m.Title != nil {
			item.Name = *m.Title
		}

		if m.Size != nil {
			item.Size = *(m.Size)
		}

		if m.Duration != nil {
			item.Duration = *(m.Duration)
		}

		if m.CreatedAt != nil {
			item.CreatedAt = *m.CreatedAt
		}

		mediaList = append(mediaList, item)
	}

	return mediaList, nil
}

func UploadVideo(ctx context.Context, publicKey string, secretKey string, data *model.UploadVideoClient, title string) error {
	client := createClient(publicKey, secretKey)
	mediaReq := aiozstreamsdk.NewCreateMediaRequest()
	mediaReq.SetTitle(title)

	media, err := client.Media.Create(*mediaReq)
	if err != nil {
		return err
	}

	mediaID := *media.GetData().Id

	errUpload := client.UploadVideo(
		ctx,
		mediaID,
		data.FileName,
		data.File,
		data.FileSize,
	)
	if errUpload != nil {
		return errUpload
	}

	return nil
}

func CreateKeyLiveStream(ctx context.Context, publicKey string, secretKey string, name string) error {
	client := createClient(publicKey, secretKey)
	typePayload := "video"
	savepayload := true

	payload := aiozstreamsdk.CreateLiveStreamKeyRequest{
		Name: &name,
		Type: &typePayload,
		Save: &savepayload,
	}
	_, err := client.LiveStream.CreateLiveStreamKeyWithContext(ctx, payload)
	if err != nil {
		return err
	}
	return nil
}
