#!/bin/bash

# Development script to run the Go Echo app with Air live reload

# Check if Air is installed
if ! command -v air &> /dev/null && ! [ -f ~/go/bin/air ]; then
    echo "❌ Air is not installed. Installing Air..."
    go install github.com/air-verse/air@latest
    if [ $? -ne 0 ]; then
        echo "❌ Failed to install Air. Please check your Go installation."
        exit 1
    fi
    echo "✅ Air installed successfully!"
fi

# Set up the proper path for Air
if command -v air &> /dev/null; then
    AIR_CMD="air"
else
    AIR_CMD="~/go/bin/air"
fi

echo "🚀 Starting development server with Air..."
echo "📁 Working directory: $(pwd)"
echo "🔄 Air will watch for file changes and auto-reload"
echo "🌐 Server will be available at http://localhost:8080"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

# Run Air
eval $AIR_CMD
