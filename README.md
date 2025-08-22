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

This section covers common continuous integration tasks that can be performed during code development. Tasks defined in this section must conform to the xc [task syntax](https://xcfile.dev/task-syntax/).

See the [xc docs](https://xcfile.dev/getting-started/) for installation instructions

### clean
Clean the cache and build artefacts
Env: BINARY_NAME=agent
Env: BUILD_DIR=bin
Env: COVERAGE_FILE=c.out
```sh
rm -rf $BUILD_DIR/$BINARY_NAME
rm -rf .agent/cache/webfetch
rm -rf $COVERAGE_FILE
```

### tidy
Run go mod tidy
```sh
go mod tidy
```

### verify
Run go mod verify
requires: tidy
```sh
go mod verify
```

### format
Format the code
```sh
go fmt  ./cmd/agent
```

### test
Test the package
```sh
go test -timeout 30s ./...
```

### test-verbose
Test the package
requires: build
```sh
go test -timeout 30s ./...
```

### test-coverage
Generate a test coverage report
requires: build
```sh
go test -coverprofile=c.out ./...
go tool cover -func=c.out
```

### dev
Run the agent in development mode. The agent will be built from the latest source.
Inputs: CLI_ARGS
Environment: CLI_ARGS=-help
```sh
go run ./cmd/agent $CLI_ARGS
```

### chat
Pass a message to the agent in development mode. The agent will be built from the latest source.
Inputs: MESSAGE
Environment: MESSAGE=hello
```sh
$MESSAGE | go run ./cmd/agent 
```

### build
Build the agent binary
requires: verify
Env: BINARY_NAME=agent
Env: BUILD_DIR=bin
```sh
mkdir -p $BUILD_DIR
go build -o $BUILD_DIR/$BINARY_NAME ./cmd/agent
```
