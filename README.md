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

## Tasks
### hello
Prints hello
```sh
echo hello
```
### world
Prints world
Requires: hello
```sh
echo world
```

### build
Build the agent binary
Env: BINARY_NAME=agent
Env: BUILD_DIR=bin
```sh
go mod tidy
mkdir -p $BUILD_DIR
go build -o $BUILD_DIR/$BINARY_NAME ./cmd/agent
```

### test
Test the package
Inputs: ARGS
```sh
go test $ARGS
```

### dev
Run the agent in development mode
Env: BINARY_NAME=agent
Env: BUILD_DIR=bin
Env: CLI_ARGS=-help
```sh
./$BUILD_DIR/$BINARY_NAME $CLI_ARGS
```