package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"

	"agent/internal/agent"
	"agent/internal/config"
)

func main() {
	// Parse command line flags
	promptFile := flag.String("f", "", "Path to prompt file for single-shot mode")
	flag.StringVar(promptFile, "prompt", "", "Path to prompt file for single-shot mode")
	continueChat := flag.Bool("continue", false, "Continue in interactive mode after processing prompt file")
	timeout := flag.Int("timeout", 60, "Timeout in seconds for non-interactive mode")
	initFlag := flag.Bool("init", false, "Initialize .agent directory")
	prePrompts := flag.String("preprompt", "", "Path to preprompt file (defaults to .agent/prompts/preprompts)")
	resumeSession := flag.String("resume", "", "Resume a specific session by ID (YYYY-MM-DD-HH-MM-SS), or use 'list' to select interactively")
	requestDelay := flag.Duration("request-delay", 0, "Delay between API requests (e.g., 2s, 500ms)")
	flag.Parse()

	// Validate flags
	if *continueChat && *promptFile == "" {
		fmt.Println("Error: --continue flag can only be used with --f/--prompt")
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Override config with command line flags
	if *requestDelay > 0 {
		cfg.Agent.RequestDelay = *requestDelay
	}

	// Create OpenAI client
	client := openai.NewClient(
		option.WithAPIKey(cfg.API.Key),
		option.WithBaseURL(cfg.API.BaseURL),
	)

	// Load pre-prompts
	prePromptList, err := getPrePrompts(*prePrompts)
	if err != nil {
		log.Fatalf("Error loading preprompts: %v", err)
	}

	// Create dependencies based on mode
	var deps agent.Dependencies
	var rl *readline.Instance

	if *promptFile != "" {
		// Single-shot mode or initial prompt with continue
		promptContent, err := loadPromptFile(*promptFile)
		if err != nil {
			log.Fatalf("Error loading prompt file: %v", err)
		}

		firstCall := true
		if *continueChat {
			rl, err = readline.New("")
			if err != nil {
				log.Fatalf("Error initializing readline: %v", err)
			}
			defer rl.Close()
		}

		deps = agent.Dependencies{
			Client: &client,
			GetUserMessage: func() (string, bool) {
				if !firstCall {
					if *continueChat {
						// Switch to interactive mode
						rl.SetPrompt("\u001b[94mYou\u001b[0m: ")
						line, err := rl.Readline()
						if err != nil {
							if err == io.EOF {
								return "", false
							}
							fmt.Printf("Error reading input: %v\n", err)
							return "", false
						}
						return line, true
					}
					return "", false
				}
				firstCall = false
				return promptContent, true
			},
			ReadlineInstance: rl,
			PrePrompts:       prePromptList,
			RequestDelay:     cfg.Agent.RequestDelay,
			SingleShot:       !*continueChat,
		}
	} else {
		// Interactive mode
		rl, err = readline.New("")
		if err != nil {
			log.Fatalf("Error initializing readline: %v", err)
		}
		defer rl.Close()

		deps = agent.Dependencies{
			Client: &client,
			GetUserMessage: func() (string, bool) {
				rl.SetPrompt("\u001b[94mYou\u001b[0m: ")
				line, err := rl.Readline()
				if err != nil {
					if err == io.EOF {
						return "", false
					}
					fmt.Printf("Error reading input: %v\n", err)
					return "", false
				}
				return line, true
			},
			ReadlineInstance: rl,
			PrePrompts:       prePromptList,
			RequestDelay:     cfg.Agent.RequestDelay,
			SingleShot:       false,
		}
	}

	// Create and configure agent
	agentApp, err := agent.New(cfg, deps)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	agentApp.SetFlags(agent.Flags{
		PromptFile:    *promptFile,
		ContinueChat:  *continueChat,
		Timeout:       time.Duration(*timeout) * time.Second,
		InitFlag:      *initFlag,
		PrePrompts:    *prePrompts,
		ResumeSession: *resumeSession,
	})

	// Set transition to interactive for continue mode
	if *continueChat {
		agentApp.SetTransitionToInteractive(true)
	}

	// Run the agent
	ctx := context.Background()
	if *promptFile != "" && !*continueChat {
		// Set timeout for single-shot mode
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(*timeout)*time.Second)
		defer cancel()
	}

	if err := agentApp.Run(ctx); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

// loadPromptFile loads and validates the prompt file
func loadPromptFile(path string) (string, error) {
	messages, err := config.LoadPromptFile(path)
	if err != nil {
		return "", err
	}
	
	// For backward compatibility with single-shot mode, combine all messages
	// TODO: Consider updating agent to handle multiple messages natively
	return strings.Join(messages, "\n\n"), nil
}

// getPrePrompts loads pre-prompts from file or returns empty slice
func getPrePrompts(prePromptsPath string) ([]string, error) {
	usingDefault := prePromptsPath == ""
	if prePromptsPath == "" {
		prePromptsPath = ".agent/prompts/preprompts"
	}

	// Check if file exists
	if _, err := os.Stat(prePromptsPath); os.IsNotExist(err) {
		if usingDefault {
			// Default file doesn't exist - that's OK, return empty
			return []string{}, nil
		} else {
			// User specified a file that doesn't exist - that's an error
			return nil, fmt.Errorf("preprompt file not found: %s", prePromptsPath)
		}
	}

	// Check file extension to determine format
	ext := strings.ToLower(filepath.Ext(prePromptsPath))
	if ext == ".yml" || ext == ".yaml" || ext == ".md" {
		// Use new LoadPromptFile function for YAML/markdown files
		messages, err := config.LoadPromptFile(prePromptsPath)
		if err != nil {
			return nil, fmt.Errorf("loading preprompts: %w", err)
		}
		return messages, nil
	}

	// Handle legacy format (plain text file with list of paths)
	content, err := os.ReadFile(prePromptsPath)
	if err != nil {
		return nil, fmt.Errorf("reading preprompts file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var prompts []string
	for _, path := range lines {
		path = strings.TrimSpace(path)
		if path != "" && !strings.HasPrefix(path, "#") {
			// Use LoadPromptFile for each individual file to support both .md and .yml
			messages, err := config.LoadPromptFile(path)
			if err != nil {
				return nil, fmt.Errorf("loading prompt file '%s': %w", path, err)
			}
			prompts = append(prompts, messages...)
		}
	}

	return prompts, nil
}
