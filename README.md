# Developing

## Setup

1. Set up your Anthropic API key:
   ```bash
   export ANTHROPIC_API_KEY="your-api-key-here"
   ```

If the 

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Lint and Test

```bash
golangci-lint fmt
golangci-lint run
```

## Build

Create a binary:
```bash
go build -o agent
```

Run the binary:
```bash
./agent
```

Or run directly without building:
```bash
go run main.go
```