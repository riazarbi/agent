# Developing

## Setup

This tool uses the OpenAI API format but can work with compatible backends. By default, it's configured for Anthropic's API, but you can customize both the API endpoint and key.

1. Configure the API endpoint:
```bash
export AGENT_BASE_URL="your-api-endpoint"  # Defaults to https://api.anthropic.com/v1/ if not set
```

2. Set up your API key using one of these environment variables:
```bash
export AGENT_API_KEY="your-api-key-here"      # Primary API key setting
# or
export ANTHROPIC_API_KEY="your-api-key-here"  # Falls back to this if AGENT_API_KEY is not set
```

2. Install dependencies:
```bash
go mod tidy
```

## Build

Create a binary:
```bash
go build -o agent
```

Run the binary in interactive mode:
```bash
./agent
```

Run with a prompt file (single-shot mode):
```bash
./agent -f prompt-file.txt
# or
./agent --prompt-file prompt-file.txt
```

Or run directly without building:
```bash
go run main.go
```