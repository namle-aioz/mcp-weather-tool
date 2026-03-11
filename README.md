# AIOZ MCP Server

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server written in Go that integrates with [AIOZ Stream](https://stream.aioz.io/) for media management and live streaming, as well as real-time weather data via the free [Open-Meteo API](https://open-meteo.com/).

---

## Features

- **Weather Lookup**: Get current temperature and wind speed for any city worldwide.
- **AIOZ Stream Integration**: Manage videos, audios, and live streams through the AIOZ Stream API.
- **Media Statistics**: Count total videos and audios in an AIOZ Stream account.
- **Video Management**: List videos, search by name, and retrieve playback URLs.
- **Video Upload**: Upload video files to AIOZ Stream directly from a Google Drive link.
- **Live Streaming**: Create live stream keys for broadcasting.

---

## MCP Tools

### `get-weather`

Get the current weather for a given location.

| Parameter  | Type   | Required | Description                    |
|------------|--------|----------|--------------------------------|
| `location` | string | Yes      | City name or location to query |

**Output example:**
```
Weather in Hanoi:
Temperature: 32.4°C
Wind Speed: 12.0 km/h
```

---

### `count-total-media`

Get the total number of videos and audios in an AIOZ Stream account.

| Parameter   | Type   | Required | Description                   |
|-------------|--------|----------|-------------------------------|
| `publicKey` | string | Yes      | AIOZ Stream public key        |
| `secretKey` | string | Yes      | AIOZ Stream secret key        |

**Output example:**
```
AIOZ Stream Account Stats:
Videos: 42
Audios: 15
```

---

### `get-video-url`

Search for a video by name and retrieve all its playback and asset URLs.

| Parameter   | Type   | Required | Description                        |
|-------------|--------|----------|------------------------------------|
| `publicKey` | string | Yes      | AIOZ Stream public key             |
| `secretKey` | string | Yes      | AIOZ Stream secret key             |
| `videoName` | string | Yes      | Name (or partial name) of the video to search |

**Output example:**
```json
{
  "EmbededURL": "https://...",
  "Mp4URL": "https://...",
  "Thumbnail": "https://...",
  "SourceURL": "https://..."
}
```

> **Note:** If multiple videos match the search term, only the last matching video's URLs are returned.

---

### `get-list-video`

Get a list of all videos in an AIOZ Stream account.

| Parameter   | Type   | Required | Description        pen    |
|-------------|--------|----------|------------------------|
| `publicKey` | string | Yes      | AIOZ Stream public key |
| `secretKey` | string | Yes      | AIOZ Stream secret key |

**Output example:**
```json
[
  {
    "MediaID": "abc123",
    "Name": "My Video",
    "Size": 104857600,
    "Duration": 120.5,
    "CreatedAt": "2024-01-15T10:00:00Z"
  }
]
```

---

### `upload-video`

Upload a video from a Google Drive link to an AIOZ Stream account. The server automatically:
1. Validates the Google Drive link format
2. Converts it to a direct download URL
3. Downloads the video file
4. Uploads it to AIOZ Stream with the specified title

| Parameter   | Type   | Required | Description                                          |
|-------------|--------|----------|------------------------------------------------------|
| `publicKey` | string | Yes      | AIOZ Stream public key                               |
| `secretKey` | string | Yes      | AIOZ Stream secret key                               |
| `videoLink` | string | Yes      | Google Drive shareable link to the video file        |
| `title`     | string | Yes      | Title/name for the uploaded video on AIOZ Stream     |

**Supported Google Drive Link Format:**
```
https://drive.google.com/file/d/{FILE_ID}/view?usp=sharing
```

**Output example:**
```
Video uploaded successfully
```

---

### `create-key-live`

Create a live stream key in an AIOZ Stream account.

| Parameter   | Type   | Required | Description                    |
|-------------|--------|----------|--------------------------------|
| `publicKey` | string | Yes      | AIOZ Stream public key         |
| `secretKey` | string | Yes      | AIOZ Stream secret key         |
| `nameKey`   | string | Yes      | Display name for the live key  |

**Output example:**
```
Key created successfully to AIOZ Stream with name: my-stream-key
```

---

## HTTP Endpoints

### `GET /ping`

Health check endpoint.

**Response:** `pong`

---

## Installation

### Prerequisites

- [Go](https://go.dev/doc/install) 1.24.0 or later
- A valid [AIOZ Stream](https://stream.aioz.io/) account with API credentials

### Build

```bash
git clone <repository-url>
cd MCP-Server-Test
go build -o aioz-mcp-server .
```

### Environment Variables

Create a `.env` file in the project root (optional):

```env
SERVER_PORT=8087
```

If `SERVER_PORT` is not set, the server defaults to port `8087`.

---

## Running the Server

```bash
go run .
```

Or using the compiled binary:

```bash
./aioz-mcp-server
```

The server will start on `http://localhost:8087` by default. MCP clients connect via SSE at `http://localhost:8087/`.

---

## Configuration

### Cursor 

To use this server with Cursor, add the following to your `mcp.json`:

```json
{
  "mcpServers": {
    "aioz-mcp-server": {
      "url": "http://localhost:your_port/sse"
    }
  }
}

```

---

## Docker

A `Dockerfile` and `docker-compose.local.yml` are provided for containerized deployment.

```bash
docker compose -f docker-compose.local.yml up --build
```

---

## Project Structure

```
.
├── main.go                        # Entry point, MCP server and tool registration
├── handler/
│   ├── aiozstream_handler.go      # MCP tool handlers for AIOZ Stream
│   └── weather_handler.go         # MCP tool handler for weather
├── tool/
│   ├── aiozstream.go              # AIOZ Stream API client logic
│   ├── weather.go                 # Open-Meteo weather API client
│   └── location.go                # Geocoding via Open-Meteo geocoding API
├── model/
│   └── client_upload_model.go     # Shared data models
├── util/
│   ├── drive.go                   # Drive func util
├── go.mod
├── Dockerfile
└── docker-compose.local.yml
```