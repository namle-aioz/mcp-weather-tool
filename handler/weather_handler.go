package handler

import (
	"context"
	"fmt"
	"mcp-weather-server/tool"

	"github.com/mark3labs/mcp-go/mcp"
)

func HandleWeather(
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

	lat, lon, err := tool.Geocode(location)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	weather, err := tool.GetWeather(lat, lon)
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
