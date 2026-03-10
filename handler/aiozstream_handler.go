package handler

import (
	"context"
	"fmt"
	"mcp-weather-server/model"
	"mcp-weather-server/tool"
	"os"
	"path/filepath"

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

func HandleUploadVideo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	fmt.Println(args)
	filePath, ok := args["filePath"].(string)
	if !ok {
		return mcp.NewToolResultError("FilePath parameter required"), nil
	}
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

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := stat.Size()
	fileName := filepath.Base(filePath)

	var clientUploadVideo = &model.UploadVideoClient{
		FileName: fileName,
		FileSize: int64(fileSize),
		File:     file,
	}

	errUpload := tool.UploadVideo(ctx, publicKey, secretKey, clientUploadVideo, videoName)
	if errUpload != nil {
		return mcp.NewToolResultError(errUpload.Error()), nil
	}

	result := fmt.Sprintf(
		"Video '%s' uploaded successfully to AIOZ Stream",
		fileName,
	)

	return mcp.NewToolResultText(result), nil
}
