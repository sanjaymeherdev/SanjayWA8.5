# WhatsApp Web Multidevice Gateway

A Go-based WhatsApp Web Multidevice gateway that exposes WhatsApp functionality over HTTP REST API and MCP (Model Context Protocol) server.

## Run & Operate

- **Run**: `bash start.sh` (builds binary then starts REST server on port 5000)
- **Build**: `cd src && CGO_ENABLED=1 go build -o whatsapp .`
- **Start REST**: `cd src && ./whatsapp rest --port=5000`
- **Start MCP**: `cd src && ./whatsapp mcp`
- **Required env vars**: None required; see `src/.env.example` for optional config

## Stack

- **Language**: Go 1.25.0
- **HTTP server**: Fiber v2 (github.com/gofiber/fiber/v2)
- **WhatsApp**: whatsmeow (go.mau.fi/whatsmeow)
- **DB**: SQLite (default) via go-sqlite3 (CGO_ENABLED=1 required), Postgres optional
- **CLI**: Cobra + Viper
- **Templates**: Fiber HTML templates (embedded)
- **MCP**: mark3labs/mcp-go over SSE

## Where things live

- `src/` — main Go module and source code
- `src/cmd/` — CLI entrypoints (rest.go, mcp.go, root.go)
- `src/domains/` — domain interfaces per feature
- `src/usecase/` — application logic
- `src/infrastructure/` — WhatsApp client + SQLite repos
- `src/ui/rest/` — HTTP REST route handlers
- `src/ui/mcp/` — MCP tool adapters
- `src/views/` — embedded HTML templates
- `src/statics/` — static assets (qrcode, senditems, media)
- `src/storages/` — SQLite database files (gitignored)
- `src/config/settings.go` — all config defaults

## Architecture decisions

- CGO is required (`go-sqlite3`) — binary must be built with `CGO_ENABLED=1`
- Embedded FS for views/assets means no separate static file deployment needed
- WhatsApp session persisted in SQLite at `storages/whatsapp.db`
- Port 5000 used (Replit webview); default upstream was 3000
- `start.sh` at repo root handles build-then-run for Replit workflow

## Product

- Scan a QR code to link a WhatsApp account
- Send messages, media, documents, stickers, polls, contacts via REST API
- Manage groups, chats, newsletters
- WebSocket endpoint for real-time events at `/ws`
- MCP server mode for AI agent integration over SSE
- Optional webhook forwarding of incoming events
- Optional basic auth protection

## User preferences

_None recorded yet._

## Gotchas

- Must build with `CGO_ENABLED=1` due to `go-sqlite3` dependency
- `ffmpeg` and `gcc` must be available at build and runtime (both present in Replit env)
- Storages directory must exist before starting; `start.sh` and `initApp()` create it automatically
- The workflow runs `start.sh` which first builds then starts — first run takes longer

## Pointers

- Upstream repo: https://github.com/aldinokemal/go-whatsapp-web-multidevice
- whatsmeow docs: https://pkg.go.dev/go.mau.fi/whatsmeow
