# passGO
An online password manager fully written in Go.

## Project Structure

```
passGO/
├── cmd/
│   ├── frontend/    # Gio desktop/web application entry point
│   │   └── main.go
│   └── backend/     # Gin HTTP server entry point
│       └── main.go
├── internal/
│   ├── frontend/    # Frontend business logic
│   │   └── app.go
│   └── backend/     # Backend business logic
│       ├── server.go
│       └── server_test.go
├── pkg/             # Shared code between frontend and backend
├── go.mod
└── go.sum
```

## Technologies

- **Frontend**: [Gio](https://gioui.org/) - Immediate mode GUI in Go
- **Backend**: [Gin](https://gin-gonic.com/) - High-performance HTTP web framework

## Getting Started

### Prerequisites

- Go 1.24 or higher
- For frontend development: X11 development libraries (Linux), or equivalent for your OS

### Building

#### Backend

```bash
go build -o passgo-backend ./cmd/backend
```

#### Frontend

```bash
go build -o passGo ./cmd/frontend
```

Configuration:

- `PASSGO_API_BASE_URL` (optional): Backend base URL used by the frontend (defaults to `http://localhost:8080`).
- `MAX_VAULTS_PER_USER` (optional): Maximum vault entries a user can create (defaults to `100`; set to `0` to disable).
- `MAX_VAULT_DATA_BYTES` (optional): Maximum plaintext vault field data per create/update request (defaults to `16384`; set to `0` to disable).

### Running

#### Backend Server

```bash
./passgo-backend
# Server starts on http://localhost:8080
```

#### Frontend Application

```bash
./passGo
```

Linux note: Gio desktop builds require system libraries (for example `libxkbcommon-dev` and Wayland/X11 dev packages) depending on your distro/window system.

### Testing

```bash
go test ./...
```

## License

MIT License - see [LICENSE](LICENSE) file for details.
