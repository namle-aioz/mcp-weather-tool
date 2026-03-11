package handler

import (
	"context"
	"fmt"
	"mcp-weather-server/model"
	"mcp-weather-server/tool"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

func HandleCountAiozStream(
	ctx context.Context,
	req mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {

	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	publicKey, ok := args["publicKey"].(string)
	if !ok {
		return mcp.NewToolResultError("PublicKey parameter required"), nil
	}
	secretKey, ok := args["secretKey"].(string)
	if !ok {
		return mcp.NewToolResultError("SecretKey parameter required"), nil
	}
	videoCount, audioCount, err := tool.CountVideoAndAudio(ctx, publicKey, secretKey)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result := fmt.Sprintf(
		"AIOZ Stream Account Stats:\nVideos: %d\nAudios: %d",
		videoCount,
		audioCount,
	)

	return mcp.NewToolResultText(result), nil
}

func HandleAiozStreamGetVideo(
	ctx context.Context,
	req mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {

	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	publicKey, ok := args["publicKey"].(string)
	if !ok {
		return mcp.NewToolResultError("PublicKey parameter required"), nil
	}
	secretKey, ok := args["secretKey"].(string)
	if !ok {
		return mcp.NewToolResultError("SecretKey parameter required"), nil
	}
	videoName, ok := args["videoName"].(string)
	if !ok {
		return mcp.NewToolResultError("SecretKey parameter required"), nil
	}

	videoDetail, err := tool.GetVideoDetailByName(ctx, publicKey, secretKey, videoName)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result := fmt.Sprintf(
		"Video Detail: %s",
		videoDetail,
	)

	return mcp.NewToolResultText(result), nil
}

func HandleAiozStreamGetListVideo(
	ctx context.Context,
	req mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {

	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	publicKey, ok := args["publicKey"].(string)
	if !ok {
		return mcp.NewToolResultError("PublicKey parameter required"), nil
	}
	secretKey, ok := args["secretKey"].(string)
	if !ok {
		return mcp.NewToolResultError("SecretKey parameter required"), nil
	}

	videos, err := tool.GetVideos(ctx, publicKey, secretKey)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result := fmt.Sprintf(
		"Video List: %s",
		videos,
	)

	return mcp.NewToolResultText(result), nil
}

func HandleUploadVideo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	fmt.Println(args)
	videoName, ok := args["title"].(string)
	if !ok {
		return mcp.NewToolResultError("title parameter required"), nil
	}
	publicKey, ok := args["publicKey"].(string)
	if !ok {
		return mcp.NewToolResultError("PublicKey parameter required"), nil
	}
	secretKey, ok := args["secretKey"].(string)
	if !ok {
		return mcp.NewToolResultError("SecretKey parameter required"), nil
	}

	uploadURL := "http://localhost:8087/upload"

	result := fmt.Sprintf(
		"Upload the video using this endpoint:\nPOST %s\nForm fields: file, title=%s, publicKey=%s, secretKey=%s",
		uploadURL,
		videoName,
		publicKey,
		secretKey,
	)

	return mcp.NewToolResultText(result), nil
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer file.Close()

	title := r.FormValue("title")
	publicKey := r.FormValue("publicKey")
	secretKey := r.FormValue("secretKey")
	if title == "" || publicKey == "" || secretKey == "" {
		http.Error(w, "title, publicKey, secretKey are required", http.StatusBadRequest)
		return
	}

	fileName := header.Filename
	fileSize := header.Size

	var clientUploadVideo = &model.UploadVideoClient{
		FileName: fileName,
		FileSize: int64(fileSize),
		File:     file,
	}

	err = tool.UploadVideo(ctx, publicKey, secretKey, clientUploadVideo, title)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte("upload success"))
}

func HandleCreateLiveStreamKey(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	nameKey, ok := args["nameKey"].(string)
	if !ok {
		return mcp.NewToolResultError("Name Key parameter required"), nil
	}
	publicKey, ok := args["publicKey"].(string)
	if !ok {
		return mcp.NewToolResultError("PublicKey parameter required"), nil
	}
	secretKey, ok := args["secretKey"].(string)
	if !ok {
		return mcp.NewToolResultError("SecretKey parameter required"), nil
	}

	err := tool.CreateKeyLiveStream(ctx, publicKey, secretKey, nameKey)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result := fmt.Sprintf(
		"Key created successfully to AIOZ Stream with name: %+v",
		nameKey,
	)

	return mcp.NewToolResultText(result), nil
}
