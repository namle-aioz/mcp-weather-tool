package main

import (
	"log"
	"mcp-weather-server/handler"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8087"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	// mux.HandleFunc("/upload", handler.UploadHandler)
	mcpServer := server.NewMCPServer(
		"aioz-mcp",
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

	aiozStream := newAIOZTool(
		"count-total-media",
		true,
		"Get total number of videos and audios in AIOZ Stream account",
	)

	aiozStreamGetVideoDetail := newAIOZTool(
		"get-video-url",
		true,
		"Get all video URL from the user's AIOZ Stream account by video name",
		mcp.WithString(
			"videoName",
			mcp.Description("Name of video"),
			mcp.Required(),
		),
	)

	aiozStreamGetListVideo := newAIOZTool(
		"get-list-video",
		true,
		"Get all video from the user's AIOZ Stream account",
	)

	aiozStreamUploadVideo := newAIOZTool(
		"upload-video",
		true,
		"Upload a video file from the user's local machine to their AIOZ Stream account",
		mcp.WithString(
			"videoLink",
			mcp.Description("URL of the video to upload (must be a Google Drive link)"),
			mcp.Required(),
		),
		mcp.WithString(
			"title",
			mcp.Description("Title of the video to upload"),
			mcp.Required(),
		),
	)

	aiozStreamUCreateKeyLiveStream := newAIOZTool(
		"create-key-live",
		true,
		"Create a key live stream to AIOZ Stream account",
		mcp.WithString(
			"nameKey",
			mcp.Description("Name of key live stream"),
			mcp.Required(),
		),
	)

	mcpServer.AddTool(weatherTool, handler.HandleWeather)
	mcpServer.AddTool(aiozStream, handler.HandleCountAiozStream)
	mcpServer.AddTool(aiozStreamGetVideoDetail, handler.HandleAiozStreamGetVideo)
	mcpServer.AddTool(aiozStreamUploadVideo, handler.HandleUploadVideo)
	mcpServer.AddTool(aiozStreamGetListVideo, handler.HandleAiozStreamGetListVideo)
	mcpServer.AddTool(aiozStreamUCreateKeyLiveStream, handler.HandleCreateLiveStreamKey)

	sseServer := server.NewSSEServer(mcpServer)
	mux.Handle("/", sseServer)
	log.Printf("Starting SSE server on port %s", port)

	http.ListenAndServe(":"+port, mux)
}

func newAIOZTool(name string, auth bool, description string, params ...mcp.ToolOption) mcp.Tool {
	baseParams := []mcp.ToolOption{
		mcp.WithDescription(description),
	}

	if auth {
		baseParams = append(baseParams,
			mcp.WithString(
				"publicKey",
				mcp.Description("Public key to authorization"),
				mcp.Required(),
			),
			mcp.WithString(
				"secretKey",
				mcp.Description("Secret key to authorization"),
				mcp.Required(),
			),
		)
	}

	baseParams = append(baseParams, params...)

	return mcp.NewTool(name, baseParams...)
}
