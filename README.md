# AIOZ MCP Server

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server written in Go that integrates with [AIOZ Stream](https://aiozstream.network/) for media management and live streaming, as well as real-time weather data via the free [Open-Meteo API](https://open-meteo.com/).

---

## Features

- **Media Statistics**: Count total videos and audios in an AIOZ Stream account.
- **Video Management**: List videos, search by name, and retrieve playback URLs.
- **Video Upload**: Upload video files to AIOZ Stream directly from a Google Drive link.
- **Live Streaming**: Create live stream keys for broadcasting.

---

## MCP Tools

---

### `count-total-media`

Get the total number of videos and audios in an AIOZ Stream account.

| Parameter   | Type   | Required | Description            |
| ----------- | ------ | -------- | ---------------------- |
| `publicKey` | string | Yes      | AIOZ Stream public key |
| `secretKey` | string | Yes      | AIOZ Stream secret key |

**Output example:**

```
AIOZ Stream Account Stats:
Videos: 42
Audios: 15
```

---

### `get-video-url`

Search for a video by name and retrieve all its playback and asset URLs.

| Parameter   | Type   | Required | Description                                   |
| ----------- | ------ | -------- | --------------------------------------------- |
| `publicKey` | string | Yes      | AIOZ Stream public key                        |
| `secretKey` | string | Yes      | AIOZ Stream secret key                        |
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

---

### `get-list-video`

Get a list of all videos in an AIOZ Stream account.

| Parameter   | Type   | Required | Description pen        |
| ----------- | ------ | -------- | ---------------------- |
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

| Parameter   | Type   | Required | Description                                      |
| ----------- | ------ | -------- | ------------------------------------------------ |
| `publicKey` | string | Yes      | AIOZ Stream public key                           |
| `secretKey` | string | Yes      | AIOZ Stream secret key                           |
| `videoLink` | string | Yes      | Google Drive shareable link to the video file    |
| `title`     | string | Yes      | Title/name for the uploaded video on AIOZ Stream |

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

| Parameter   | Type   | Required | Description                   |
| ----------- | ------ | -------- | ----------------------------- |
| `publicKey` | string | Yes      | AIOZ Stream public key        |
| `secretKey` | string | Yes      | AIOZ Stream secret key        |
| `nameKey`   | string | Yes      | Display name for the live key |

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

> **Trước khi bắt đầu:** Đảm bảo server đang chạy bằng lệnh `go run .` trong thư mục project. Bạn sẽ thấy dòng log `Starting SSE server on port 8087`.

---

### 🖱️ Cursor

Cursor là code editor hỗ trợ MCP natively. Làm theo các bước sau:

**Bước 1 — Mở file cấu hình MCP**

- Nhấn `Ctrl + Shift + P` (Windows/Linux) hoặc `Cmd + Shift + P` (macOS)
- Gõ **"Open MCP Config"** → nhấn Enter
- File `mcp.json` sẽ mở ra (thường nằm tại `~/.cursor/mcp.json`)

**Bước 2 — Thêm server vào config**

Dán nội dung sau vào file (nếu file đã có sẵn nội dung, chỉ thêm phần `aioz-mcp-server` vào trong `mcpServers`):

```json
{
  "mcpServers": {
    "aioz-mcp-server": {
      "url": "http://localhost:8087/sse"
    }
  }
}
```

> 💡 Thay `8087` bằng port bạn đã cấu hình trong file `.env` nếu khác.

**Bước 3 — Lưu file và kiểm tra**

- Lưu file (`Ctrl + S`)
- Mở **Cursor Chat** (`Ctrl + L`) và hỏi thử: _"Đếm tổng số video trong tài khoản AIOZ Stream của tôi"_
- Cursor sẽ tự động gọi tool `count-total-media`

![Cursor MCP Config](docs/images/cursor-mcp-config.png)

---

### 🤖 Claude Desktop

Claude Desktop hỗ trợ MCP thông qua file config riêng.

**Bước 1 — Build server thành file binary**

```bash
go build -o aioz-mcp-server .
```

> Sau lệnh này sẽ có file `aioz-mcp-server.exe` (Windows) hoặc `aioz-mcp-server` (macOS/Linux) trong thư mục project.

**Bước 2 — Mở file cấu hình Claude Desktop**

Mở file tại đường dẫn sau (tạo mới nếu chưa có):

| Hệ điều hành | Đường dẫn                                                         |
| ------------ | ----------------------------------------------------------------- |
| **Windows**  | `%APPDATA%\Claude\claude_desktop_config.json`                     |
| **macOS**    | `~/Library/Application Support/Claude/claude_desktop_config.json` |
| **Linux**    | `~/.config/Claude/claude_desktop_config.json`                     |

**Bước 3 — Thêm server vào config**

```json
{
  "mcpServers": {
    "aioz-mcp-server": {
      "command": "C:/path/to/aioz-mcp-server.exe",
      "args": [],
      "env": {
        "SERVER_PORT": "8087"
      }
    }
  }
}
```

> ⚠️ Thay `C:/path/to/aioz-mcp-server.exe` bằng **đường dẫn tuyệt đối** đến file binary bạn vừa build ở Bước 1.  
> **Ví dụ Windows:** `C:/Users/YourName/projects/mcp-weather-tool/aioz-mcp-server.exe`

**Bước 4 — Khởi động lại Claude Desktop**

- Đóng hoàn toàn Claude Desktop (kể cả system tray)
- Mở lại và kiểm tra: gõ _"Liệt kê tất cả video trong AIOZ Stream"_ → Claude sẽ gọi tool `get-list-video`

![Claude Desktop MCP Config](docs/images/claude-mcp-config.png)

---

### 💬 ChatGPT Desktop

ChatGPT Desktop (macOS) hỗ trợ MCP từ phiên bản mới nhất.

**Bước 1 — Cài đặt ChatGPT Desktop**

Tải về tại: [https://openai.com/chatgpt/download](https://openai.com/chatgpt/download)

**Bước 2 — Mở Settings → Developers**

- Click vào avatar góc trên bên trái
- Chọn **Settings** → **Developers**
- Bật **"Enable MCP Servers"** (nếu có toggle)

**Bước 3 — Thêm MCP Server**

- Click **"Add Server"** hoặc **"Edit Config"**
- Chọn loại kết nối là **SSE**
- Nhập URL: `http://localhost:8087/sse`

**Bước 4 — Kiểm tra**

Trong cửa sổ chat, hỏi: _"Get me the weather in Ho Chi Minh City"_  
ChatGPT sẽ tự động gọi tool `get-weather`.

![ChatGPT MCP Config](docs/images/chatgpt-mcp-config.png)

### 🔵 Visual Studio Code (VS Code)

VS Code hỗ trợ MCP thông qua **GitHub Copilot** (cần bật Agent mode).

> ⚠️ Tính năng này yêu cầu VS Code phiên bản **1.99+** và có đăng ký **GitHub Copilot**.

**Bước 1 — Mở file cấu hình**

Có 2 cách:

- **Cách A (workspace):** Tạo file `.vscode/mcp.json` trong thư mục project của bạn _(chỉ áp dụng cho project đó)_
- **Cách B (toàn máy):** Mở `settings.json` → nhấn `Ctrl + Shift + P` → gõ `Open User Settings JSON`

**Bước 2 — Thêm cấu hình MCP**

Nếu dùng **Cách A** — tạo file `.vscode/mcp.json`:

```json
{
  "servers": {
    "aioz-mcp-server": {
      "type": "sse",
      "url": "http://localhost:8087/sse"
    }
  }
}
```

Nếu dùng **Cách B** — thêm vào `settings.json`:

```json
{
  "mcp": {
    "servers": {
      "aioz-mcp-server": {
        "type": "sse",
        "url": "http://localhost:8087/sse"
      }
    }
  }
}
```

> 💡 Thay `8087` bằng port bạn đã cấu hình nếu khác.

**Bước 3 — Bật Agent Mode trong Copilot Chat**

- Mở **Copilot Chat** (`Ctrl + Alt + I`)
- Ở góc trên của chat, chuyển sang chế độ **Agent** (thay vì Ask/Edit)
- Gõ thử: _"Lấy danh sách video trong AIOZ Stream của tôi"_

![VS Code MCP Config](docs/images/vscode-mcp-config.png)

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
├── tool/
│   ├── aiozstream.go              # AIOZ Stream API client logic
├── model/
│   └── client_upload_model.go     # Shared data models
├── util/
│   ├── drive.go                   # Drive func util
├── go.mod
├── Dockerfile
└── docker-compose.local.yml
```
