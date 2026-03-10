package main

import (
	"log"
	"mcp-weather-server/handler"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	mcpServer := server.NewMCPServer(
		"weather-server",
		"1.0.0",
	)

	weatherTool := mcp.NewTool(
		"get-weather",
		mcp.WithDescription("Get the current weather for a location"),
		mcp.WithString(
			"location",
			mcp.Description("City name or location"),
			mcp.Required(),
		),
	)

	aiozStream := mcp.NewTool(
		"count-total-media",
		mcp.WithDescription("Get total number of videos and audios in AIOZ Stream account"),
		mcp.WithString(
			"publicKey",
			mcp.Description("Public key to authorization"),
			mcp.Required(),
		), mcp.WithString(
			"secretKey",
			mcp.Description("Secret key to authorization"),
			mcp.Required(),
		),
	)

	aiozStreamGetVideoDetail := mcp.NewTool(
		"get-video-url",
		mcp.WithDescription("Get all video URL from the user's AIOZ Stream account by video name"),
		mcp.WithString(
			"publicKey",
			mcp.Description("Public key to authorization"),
			mcp.Required(),
		), mcp.WithString(
			"secretKey",
			mcp.Description("Secret key to authorization"),
			mcp.Required(),
		), mcp.WithString(
			"videoName",
			mcp.Description("Name of video"),
			mcp.Required(),
		),
	)

	aiozStreamUploadVideo := mcp.NewTool(
		"upload-video",
		mcp.WithDescription("Upload a video file from the user's local machine to their AIOZ Stream account using the provided file path "),
		mcp.WithString(
			"publicKey",
			mcp.Description("Public key to authorization"),
			mcp.Required(),
		), mcp.WithString(
			"secretKey",
			mcp.Description("Secret key to authorization"),
			mcp.Required(),
		), mcp.WithString(
			"filePath",
			mcp.Description("Path of this file in local machine"),
			mcp.Required(),
		), mcp.WithString(
			"title",
			mcp.Description("Title of the video to upload"),
			mcp.Required(),
		),
	)

	mcpServer.AddTool(weatherTool, handler.HandleWeather)
	mcpServer.AddTool(aiozStream, handler.HandleCountAiozStream)
	mcpServer.AddTool(aiozStreamGetVideoDetail, handler.HandleAiozStreamGetVideo)
	mcpServer.AddTool(aiozStreamUploadVideo, handler.HandleUploadVideo)

	sseServer := server.NewSSEServer(mcpServer)
	log.Printf("Starting SSE server on localhost:8087")

	if err := sseServer.Start(":8087"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
