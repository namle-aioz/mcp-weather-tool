package main

import (
	"context"
	"fmt"
	"log"

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
			"location",
			mcp.Description("City name or location"),
			mcp.Required(),
		),
	)

	mcpServer.AddTool(weatherTool, handleWeather)
	mcpServer.AddTool(aiozStream, handleAiozStream)

	sseServer := server.NewSSEServer(mcpServer)
	log.Printf("Starting SSE server on localhost:8087")

	if err := sseServer.Start(":8087"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func handleWeather(
	ctx context.Context,
	req mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {

	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}
	location, ok := args["location"].(string)
	if !ok {
		return mcp.NewToolResultError("location parameter required"), nil
	}

	lat, lon, err := geocode(location)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	weather, err := getWeather(lat, lon)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result := fmt.Sprintf(
		"Weather in %s:\nTemperature: %.1f°C\nWind Speed: %.1f km/h",
		location,
		weather.CurrentWeather.Temperature,
		weather.CurrentWeather.WindSpeed,
	)

	return mcp.NewToolResultText(result), nil

}
func handleAiozStream(
	ctx context.Context,
	req mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {

	videoCount, audioCount, err := countVideoAndAudio(ctx)
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
