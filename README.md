# Go Echo Web App

A simple web application built with Go and the Echo framework.

## Features

- RESTful API endpoints
- JSON responses
- Middleware for logging, recovery, and CORS
- Basic user management endpoints
- Health check endpoint
- **TTS (Text-to-Speech) API** with Volcengine integration
- Basic audio file management
- Environment-based configuration

## Prerequisites

- Go 1.24.6 or higher
- Git
- Air (for live reload during development)

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd common-go
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Configure TTS service (optional):
   ```bash
   # Copy the example environment file
   cp .env.example .env

   # Edit .env file with your Volcengine TTS credentials
   # Get credentials from: https://console.volcengine.com/speech/app
   ```

4. Install Air for live reload (optional but recommended for development):
   ```bash
   go install github.com/air-verse/air@latest
   ```

   Note: Make sure `~/go/bin` is in your PATH, or use the full path `~/go/bin/air`

## Running the Application

### Quick Start (Recommended)
Use the provided development script or Makefile:
```bash
# Using the development script
./dev.sh

# Or using Make
make dev
```

### Development Mode (with live reload)
Start the server with Air for automatic reloading on file changes:
```bash
# If ~/go/bin is in your PATH
air

# Or use full path
~/go/bin/air
```

### Production Mode
Start the server normally:
```bash
# Direct run
go run main.go

# Or build first, then run
make build && ./app

# Or build and run in one command
make run
```

The server will start on `http://localhost:8080`

## API Endpoints

### General
- `GET /` - Welcome message
- `GET /health` - Health check

### TTS (Text-to-Speech) API
- `GET /api/tts?text=your_text` - Convert text to speech
- `GET /api/tts/files` - List all audio files
- `GET /api/tts/status` - Get TTS service status (when disabled)

## Example Usage

### Get all users
```bash
curl http://localhost:8080/api/users
```

### Create a user
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Johnson", "email": "alice@example.com"}'
```

### Get a specific user
```bash
curl http://localhost:8080/api/users/1
```

### Health check
```bash
curl http://localhost:8080/health
```

### TTS Examples

#### Synthesize speech
```bash
curl "http://localhost:8080/api/tts?text=Hello,%20this%20is%20a%20test"
```

#### List audio files
```bash
curl http://localhost:8080/api/tts/files
```

#### Check TTS status (when disabled)
```bash
curl http://localhost:8080/api/tts/status
```

### Complete TTS Examples



## Project Structure

```
common-go/
├── main.go          # Main application file
├── go.mod           # Go module file
├── go.sum           # Go dependencies checksum
├── .air.toml        # Air configuration file
├── .env.example     # Environment variables example
├── .gitignore       # Git ignore file
├── Makefile         # Common development tasks
├── dev.sh           # Development script with Air
├── tmp/             # Temporary build files (created by Air)
├── audio/           # Generated audio files (created automatically)
├── tts/             # TTS package
│   ├── client.go    # Volcengine TTS client
│   ├── handlers.go  # HTTP handlers for TTS endpoints
│   ├── config.go    # Configuration management
│   └── model.go     # Data models and types
└── README.md        # This file
```

## Built With

- [Go](https://golang.org/) - Programming language
- [Echo](https://echo.labstack.com/) - High performance, minimalist Go web framework
- [Volcengine TTS](https://www.volcengine.com/docs/6561/1257544) - Text-to-Speech API
- [Air](https://github.com/air-verse/air) - Live reload for Go apps
- [GoDotEnv](https://github.com/joho/godotenv) - Environment variable loader

## Development

### Available Make Commands

```bash
make help        # Show all available commands
make dev         # Start development server with Air (live reload)
make build       # Build the application
make run         # Build and run the application
make clean       # Clean build artifacts and temporary files
make test        # Run tests
make install-air # Install Air for live reload
make deps        # Download and tidy up dependencies
make fmt         # Format Go code
make vet         # Run go vet
make check       # Run all checks (format, vet, test)
make test-api    # Test all API endpoints
make setup       # Initial project setup
```

### Development Workflow

1. **Initial setup**: `make setup`
2. **Start developing**: `make dev` (starts Air with live reload)
3. **Before committing**: `make check` (formats, vets, and tests code)
4. **Test API endpoints**: `make test-api` (tests all endpoints)
5. **Build for production**: `make build`

To add new routes, modify the `main.go` file and add your handlers. The application uses Echo's routing system with middleware for common functionality like logging and CORS.

### Live Reload with Air

Air is configured to:
- Watch for changes in `.go` files
- Automatically rebuild and restart the server
- Exclude test files and temporary directories
- Log build errors to `build-errors.log`

You can customize Air's behavior by editing the `.air.toml` configuration file.

### TTS Configuration

To use the TTS service, you need to:

1. **Get Volcengine Credentials**: Visit [Volcengine Console](https://console.volcengine.com/speech/app) and create a TTS application
2. **Configure Environment**: Copy `.env.example` to `.env` and fill in your credentials:
   ```bash
   VOLCENGINE_TTS_APP_ID=your_app_id
   VOLCENGINE_TTS_ACCESS_KEY=your_access_key
   VOLCENGINE_TTS_SECRET_KEY=your_secret_key
   TTS_ENABLED=true
   ```
3. **Audio Storage**: Audio files are automatically saved to `./audio/` directory
4. **Usage**: Send text via query parameter to `/api/tts?text=your_text`

#### TTS API Response Format

```json
{
  "file_path": "./audio/tts_your_text.mp3",
  "code": 0
}
```



The TTS service uses default voice settings from the Volcengine configuration.

For production deployment, consider:
- Adding environment configuration
- Implementing proper database integration
- Adding authentication and authorization
- Setting up proper error handling
- Adding input validation
- Implementing logging to files
- Adding graceful shutdown handling
- Securing TTS API with rate limiting
- Adding audio file cleanup policies

## License

This project is open source and available under the MIT License.
