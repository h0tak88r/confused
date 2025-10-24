package main

import (
	"fmt"
	"os"

	"github.com/h0tak88r/confused2/pkg/config"
	"github.com/h0tak88r/confused2/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	cfg    *config.Config
	log    *logger.Logger
	version = "2.2.0"
	buildDate = "2025-10-24"
)

func main() {
	var err error
	
	// Initialize config
	cfg = config.Default()
	
	// Initialize logger
	log, err = logger.New(logger.INFO, cfg.Verbose, "")
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	// Setup CLI
	rootCmd := setupRootCommand()
	
	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Error("Command execution failed: %v", err)
		os.Exit(1)
	}
}

func setupRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "confused",
		Short: "Advanced Dependency Confusion Scanner",
		Long: `Confused is an advanced dependency confusion scanner that can:
- Scan local dependency files
- Scan GitHub repositories and organizations
- Discover dependency files via web scanning
- Support multiple package managers concurrently
- Generate comprehensive reports`,
		Version: fmt.Sprintf("%s (built %s)", version, buildDate),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Load configuration
			config.Load(cfg)
			
			// Setup logger with new settings
			log.SetVerbose(cfg.Verbose)
			if cfg.Verbose {
				log.SetLevel(logger.DEBUG)
			}
		},
	}

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringVarP(&cfg.Output, "output", "o", "", "Output file path")
	rootCmd.PersistentFlags().StringVarP(&cfg.Format, "format", "f", "text", "Output format (text, json, html)")
	rootCmd.PersistentFlags().IntVarP(&cfg.Workers, "workers", "w", 10, "Number of concurrent workers")
	rootCmd.PersistentFlags().IntVar(&cfg.Timeout, "timeout", 30, "Request timeout in seconds")
	rootCmd.PersistentFlags().StringSliceVar(&cfg.SafeSpaces, "safe-spaces", []string{}, "Known-safe namespaces (supports wildcards)")
	rootCmd.PersistentFlags().StringVar(&cfg.OutputDir, "output-dir", "./results", "Output directory for results")
	rootCmd.PersistentFlags().BoolVar(&cfg.SaveResults, "save", true, "Save results to files")
	rootCmd.PersistentFlags().StringVar(&cfg.GitHubToken, "github-token", "", "GitHub API token")
	rootCmd.PersistentFlags().StringVar(&cfg.UserAgent, "user-agent", "Confused-DepConfusion-Scanner/2.0", "User agent for HTTP requests")

	// Add subcommands
	rootCmd.AddCommand(createScanCommand())
	rootCmd.AddCommand(createGitHubCommand())
	rootCmd.AddCommand(createWebCommand())
	rootCmd.AddCommand(createConfigCommand())

	return rootCmd
}
