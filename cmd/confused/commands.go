package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/h0tak88r/confused/internal/resolvers"
	"github.com/h0tak88r/confused/internal/types"
	"github.com/h0tak88r/confused/pkg/github"
	"github.com/h0tak88r/confused/pkg/web"
	"github.com/spf13/cobra"
)

// createScanCommand creates the scan command
func createScanCommand() *cobra.Command {
	scanCmd := &cobra.Command{
		Use:   "scan [flags] <file>",
		Short: "Scan a local dependency file for dependency confusion vulnerabilities",
		Long: `Scan a local dependency file for dependency confusion vulnerabilities.
Supports multiple package managers: npm, pip, composer, mvn, rubygems.

Examples:
  confused scan package.json
  confused scan -l pip requirements.txt
  confused scan -l mvn pom.xml -w 20
  confused scan -l composer composer.json --safe-spaces "@mycompany/*"`,
		Args: cobra.ExactArgs(1),
		RunE: runScanCommand,
	}

	scanCmd.Flags().StringP("language", "l", "npm", "Package manager (npm, pip, composer, mvn, rubygems)")
	scanCmd.Flags().StringSlice("safe-spaces", []string{}, "Known-safe namespaces (supports wildcards)")

	return scanCmd
}

// createGitHubCommand creates the GitHub scanning command
func createGitHubCommand() *cobra.Command {
	githubCmd := &cobra.Command{
		Use:   "github",
		Short: "GitHub scanning commands",
		Long:  `Scan GitHub repositories and organizations for dependency confusion vulnerabilities.`,
	}

	// GitHub repo scan
	repoCmd := &cobra.Command{
		Use:   "repo <owner/repo>",
		Short: "Scan a specific GitHub repository",
		Long: `Scan a specific GitHub repository for dependency files and check for dependency confusion vulnerabilities.

Examples:
  confused github repo microsoft/PowerShell
  confused github repo facebook/react -w 20
  confused github repo google/go --languages npm,pip`,
		Args: cobra.ExactArgs(1),
		RunE: runGitHubRepoCommand,
	}

	repoCmd.Flags().StringSlice("languages", []string{"npm", "pip", "composer", "mvn", "rubygems"}, "Package managers to scan for")
	repoCmd.Flags().StringSlice("safe-spaces", []string{}, "Known-safe namespaces (supports wildcards)")
	repoCmd.Flags().Bool("deep", false, "Perform deep scan including all branches")
	repoCmd.Flags().IntP("workers", "w", 10, "Number of concurrent workers")

	// GitHub org scan
	orgCmd := &cobra.Command{
		Use:   "org <organization>",
		Short: "Scan a GitHub organization",
		Long: `Scan all repositories in a GitHub organization for dependency confusion vulnerabilities.

Examples:
  confused github org microsoft
  confused github org google --max-repos 100
  confused github org facebook --languages npm,pip --workers 20`,
		Args: cobra.ExactArgs(1),
		RunE: runGitHubOrgCommand,
	}

	orgCmd.Flags().StringSlice("languages", []string{"npm", "pip", "composer", "mvn", "rubygems"}, "Package managers to scan for")
	orgCmd.Flags().StringSlice("safe-spaces", []string{}, "Known-safe namespaces (supports wildcards)")
	orgCmd.Flags().Int("max-repos", 50, "Maximum number of repositories to scan")
	orgCmd.Flags().Bool("deep", false, "Perform deep scan including all branches")
	orgCmd.Flags().IntP("workers", "w", 10, "Number of concurrent workers")

	githubCmd.AddCommand(repoCmd)
	githubCmd.AddCommand(orgCmd)

	return githubCmd
}

// createWebCommand creates the web scanning command
func createWebCommand() *cobra.Command {
	webCmd := &cobra.Command{
		Use:   "web <target> [target2] [target3] ...",
		Short: "Scan web targets for dependency files",
		Long: `Scan web targets to discover and analyze dependency files for dependency confusion vulnerabilities.
This command will attempt to discover dependency files through various methods including:
- Common file paths (package.json, requirements.txt, etc.)
- Directory brute forcing
- Sitemap analysis
- Response analysis

Examples:
  confused web https://example.com
  confused web https://example.com https://app.example.com
  confused web example.com --deep --workers 20
  confused web --target-file targets.txt`,
		Args: func(cmd *cobra.Command, args []string) error {
			targetFile, _ := cmd.Flags().GetString("target-file")
			if targetFile != "" {
				// If target file is provided, no command line args are required
				return nil
			}
			// Otherwise, require at least one target
			if len(args) < 1 {
				return fmt.Errorf("requires at least 1 arg(s), only received %d", len(args))
			}
			return nil
		},
		RunE: runWebCommand,
	}

	webCmd.Flags().StringSlice("languages", []string{"npm", "pip", "composer", "mvn", "rubygems"}, "Package managers to scan for")
	webCmd.Flags().StringSlice("safe-spaces", []string{}, "Known-safe namespaces (supports wildcards)")
	webCmd.Flags().Bool("deep", false, "Perform deep scan with extensive file discovery")
	webCmd.Flags().StringSlice("wordlist", []string{}, "Custom wordlist for file discovery")
	webCmd.Flags().Int("max-depth", 3, "Maximum directory depth for discovery")
	webCmd.Flags().String("target-file", "", "File containing list of targets (one per line)")
	webCmd.Flags().IntP("workers", "w", 10, "Number of concurrent workers")

	return webCmd
}

// createConfigCommand creates the config command
func createConfigCommand() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management",
		Long:  `Manage configuration settings for the confused scanner.`,
	}

	// Generate config
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a sample configuration file",
		Long:  `Generate a sample configuration file with default settings.`,
		RunE:  runConfigGenerateCommand,
	}

	generateCmd.Flags().StringP("output", "o", "confused.yaml", "Output configuration file path")

	// Validate config
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration file",
		Long:  `Validate the current configuration file for errors.`,
		RunE:  runConfigValidateCommand,
	}

	configCmd.AddCommand(generateCmd)
	configCmd.AddCommand(validateCmd)

	return configCmd
}

// runScanCommand runs the scan command
func runScanCommand(cmd *cobra.Command, args []string) error {
	filename := args[0]
	language, _ := cmd.Flags().GetString("language")
	safeSpaces, _ := cmd.Flags().GetStringSlice("safe-spaces")

	log.Info("Starting dependency confusion scan...")
	log.Info("Target: %s", filename)
	log.Info("Language: %s", language)
	log.Info("Workers: %d", cfg.Workers)

	// Create scan result
	result := types.NewScanResult(filename, "file", language)

	// Get resolver for the language
	resolver, err := resolvers.GetResolverForLanguageWithVerbose(language, cfg.Verbose)
	if err != nil {
		return fmt.Errorf("failed to get resolver for language %s: %w", language, err)
	}

	// Read packages from file
	if err := resolver.ReadPackagesFromFile(filename); err != nil {
		return fmt.Errorf("failed to read packages from file: %w", err)
	}

	// Get vulnerable packages
	vulnerablePackages := resolver.PackagesNotInPublic()
	
	// Remove safe spaces
	vulnerablePackages = removeSafe(vulnerablePackages, safeSpaces)

	// Add to result
	for _, pkg := range vulnerablePackages {
		result.AddVulnerable(pkg)
	}

	// Finalize result
	result.Finalize()

	// Print result
	printResult(result)

	// Save results if requested
	if cfg.SaveResults {
		if err := saveResults(result); err != nil {
			log.Warn("Failed to save results: %v", err)
		}
	}

	// Exit with error code if vulnerable
	if result.IsVulnerable() {
		os.Exit(1)
	}

	return nil
}

// runGitHubRepoCommand runs the GitHub repository scan command
func runGitHubRepoCommand(cmd *cobra.Command, args []string) error {
	repo := args[0]
	languages, _ := cmd.Flags().GetStringSlice("languages")
	safeSpaces, _ := cmd.Flags().GetStringSlice("safe-spaces")
	deep, _ := cmd.Flags().GetBool("deep")
	workers, _ := cmd.Flags().GetInt("workers")

	log.Info("Starting GitHub repository scan...")
	log.Info("Repository: %s", repo)
	log.Info("Languages: %s", strings.Join(languages, ", "))
	log.Info("Deep scan: %v", deep)
	log.Info("Workers: %d", workers)

	// Check if GitHub token is provided
	if cfg.GitHubToken == "" {
		log.Warn("No GitHub token provided. Using unauthenticated requests (rate limited)")
		log.Warn("For better performance and higher rate limits, use --github-token flag")
	}

	// Initialize GitHub client
	githubClient, err := github.New(cfg, log)
	if err != nil {
		return fmt.Errorf("failed to initialize GitHub client: %w", err)
	}

	// Scan repository
	results, err := githubClient.ScanRepository(repo, languages, safeSpaces, deep)
	if err != nil {
		return fmt.Errorf("failed to scan repository: %w", err)
	}

	// Process results
	for _, result := range results {
		printResult(result)
	}

	// Save results if requested
	if cfg.SaveResults {
		// Convert []*types.ScanResult to []types.ScanResult
		convertedResults := make([]types.ScanResult, len(results))
		for i, result := range results {
			convertedResults[i] = *result
		}
		if err := saveScanResults(convertedResults); err != nil {
			log.Warn("Failed to save results: %v", err)
		}
	}

	return nil
}

// runGitHubOrgCommand runs the GitHub organization scan command
func runGitHubOrgCommand(cmd *cobra.Command, args []string) error {
	org := args[0]
	languages, _ := cmd.Flags().GetStringSlice("languages")
	safeSpaces, _ := cmd.Flags().GetStringSlice("safe-spaces")
	maxRepos, _ := cmd.Flags().GetInt("max-repos")
	deep, _ := cmd.Flags().GetBool("deep")
	workers, _ := cmd.Flags().GetInt("workers")

	log.Info("Starting GitHub organization scan...")
	log.Info("Organization: %s", org)
	log.Info("Languages: %s", strings.Join(languages, ", "))
	log.Info("Max repositories: %d", maxRepos)
	log.Info("Deep scan: %v", deep)
	log.Info("Workers: %d", workers)

	// Check if GitHub token is provided
	if cfg.GitHubToken == "" {
		log.Warn("No GitHub token provided. Using unauthenticated requests (rate limited)")
		log.Warn("For better performance and higher rate limits, use --github-token flag")
	}

	// Initialize GitHub client
	githubClient, err := github.New(cfg, log)
	if err != nil {
		return fmt.Errorf("failed to initialize GitHub client: %w", err)
	}

	// Scan organization
	results, err := githubClient.ScanOrganization(org, languages, safeSpaces, maxRepos, deep)
	if err != nil {
		return fmt.Errorf("failed to scan organization: %w", err)
	}

	// Process results
	for _, result := range results {
		printResult(result)
	}

	// Save results if requested
	if cfg.SaveResults {
		// Convert []*types.ScanResult to []types.ScanResult
		convertedResults := make([]types.ScanResult, len(results))
		for i, result := range results {
			convertedResults[i] = *result
		}
		if err := saveScanResults(convertedResults); err != nil {
			log.Warn("Failed to save results: %v", err)
		}
	}

	return nil
}

// runWebCommand runs the web scanning command
func runWebCommand(cmd *cobra.Command, args []string) error {
	languages, _ := cmd.Flags().GetStringSlice("languages")
	deep, _ := cmd.Flags().GetBool("deep")
	maxDepth, _ := cmd.Flags().GetInt("max-depth")
	targetFile, _ := cmd.Flags().GetString("target-file")
	workers, _ := cmd.Flags().GetInt("workers")

	// Get targets from command line args or target file
	var targets []string
	if targetFile != "" {
		fileTargets, err := readTargetFile(targetFile)
		if err != nil {
			return fmt.Errorf("failed to read target file: %w", err)
		}
		targets = fileTargets
	} else {
		targets = args
	}

	log.Info("Starting web target scan...")
	log.Info("Targets: %d", len(targets))
	log.Info("Languages: %s", strings.Join(languages, ", "))
	log.Info("Deep scan: %v", deep)
	log.Info("Max depth: %d", maxDepth)
	log.Info("Workers: %d", workers)

	// Initialize web scanner
	webScanner := web.New(log, cfg.UserAgent, cfg.GetTimeout())

	// Process targets with worker pool
	results, err := processWebTargetsWithWorkers(webScanner, targets, languages, deep, maxDepth, workers)
	if err != nil {
		return fmt.Errorf("failed to scan web targets: %w", err)
	}

	// Process results
	for _, result := range results {
		printResult(result)
	}

	// Save results if requested
	if cfg.SaveResults {
		// Convert []*types.ScanResult to []types.ScanResult
		convertedResults := make([]types.ScanResult, len(results))
		for i, result := range results {
			convertedResults[i] = *result
		}
		if err := saveScanResults(convertedResults); err != nil {
			log.Warn("Failed to save results: %v", err)
		}
	}

	return nil
}

// runConfigGenerateCommand runs the config generate command
func runConfigGenerateCommand(cmd *cobra.Command, args []string) error {
	output, _ := cmd.Flags().GetString("output")

	log.Info("Generating sample configuration file...")

	// Create sample config
	sampleConfig := `# Confused Dependency Confusion Scanner Configuration

# General settings
verbose: false
output: ""
format: "text"  # text, json, html
workers: 10
timeout: 30

# GitHub settings
github_token: ""  # Set via CONFUSED_GITHUB_TOKEN environment variable or --github-token flag
github_org: ""
github_repo: ""
max_repos: 50

# Target settings
targets: []
target_file: ""

# Scanning settings
safe_spaces: []  # Known-safe namespaces (supports wildcards)
languages: ["npm", "pip", "composer", "mvn", "rubygems"]
deep_scan: false

# Rate limiting
rate_limit: 100
delay: 100

# Output settings
save_results: true
output_dir: "./results"

# Web scanning
user_agent: "Confused-DepConfusion-Scanner/2.0"
follow_redirects: true

# Database settings (for AutoAR integration)
database:
  type: "postgresql"  # postgresql, sqlite
  host: "localhost"
  port: 5432
  username: ""
  password: ""
  database: "confused"
`

	// Write config file
	if err := os.WriteFile(output, []byte(sampleConfig), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	log.Info("Sample configuration file generated: %s", output)
	return nil
}

// runConfigValidateCommand runs the config validate command
func runConfigValidateCommand(cmd *cobra.Command, args []string) error {
	log.Info("Validating configuration...")

	// Check if config is valid
	if cfg == nil {
		return fmt.Errorf("configuration is nil")
	}

	// Validate required fields
	if cfg.Workers <= 0 {
		return fmt.Errorf("workers must be greater than 0")
	}

	if cfg.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	// Validate languages
	validLanguages := []string{"npm", "pip", "composer", "mvn", "rubygems"}
	for _, lang := range cfg.Languages {
		valid := false
		for _, validLang := range validLanguages {
			if lang == validLang {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid language: %s", lang)
		}
	}

	log.Info("Configuration is valid")
	return nil
}

// printResult outputs the result of the scanner
func printResult(result *types.ScanResult) {
	if !result.IsVulnerable() {
		log.Info("All packages seem to be available in the public repositories.")
		log.Info("In case your application uses private repositories please make sure that those namespaces in public repositories are controlled by a trusted party.")
		return
	}
	
	log.Warn("Issues found, the following packages are not available in public package repositories:")
	for _, pkg := range result.Vulnerable {
		log.Warn(" [!] %s", pkg)
	}
}

// saveResults saves a single scan result
func saveResults(result *types.ScanResult) error {
	// Ensure output directory exists
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("confused-scan-%s-%s.json", result.Target, timestamp)
	filepath := filepath.Join(cfg.OutputDir, filename)

	// Marshal result
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	// Write file
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write result file: %w", err)
	}

	log.Info("Results saved to: %s", filepath)
	return nil
}

// saveScanResults saves multiple scan results
func saveScanResults(results []types.ScanResult) error {
	// Ensure output directory exists
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("confused-results-%s.json", timestamp)
	filepath := filepath.Join(cfg.OutputDir, filename)

	// Create scan results
	scanResults := &types.ScanResults{
		Results: results,
	}

	// Calculate summary
	scanResults.Summary.TotalTargets = len(results)
	for _, result := range results {
		scanResults.Summary.VulnerableCount += len(result.Vulnerable)
		scanResults.Summary.SafeCount += len(result.Safe)
		if result.Duration > scanResults.Summary.Duration {
			scanResults.Summary.Duration = result.Duration
		}
	}

	// Marshal results
	data, err := json.MarshalIndent(scanResults, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	// Write file
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write results file: %w", err)
	}

	log.Info("Results saved to: %s", filepath)
	return nil
}

// removeSafe removes known-safe package names from the slice
func removeSafe(packages []string, safeSpaces []string) []string {
	if len(safeSpaces) == 0 {
		return packages
	}
	
	retSlice := []string{}
	for _, pkg := range packages {
		ignored := false
		for _, safeSpace := range safeSpaces {
			ok, err := filepath.Match(safeSpace, pkg)
			if err != nil {
				log.Warn("Encountered an error while trying to match a known-safe namespace %s: %v", safeSpace, err)
				continue
			}
			if ok {
				ignored = true
				break
			}
		}
		if !ignored {
			retSlice = append(retSlice, pkg)
		}
	}
	return retSlice
}

// readTargetFile reads targets from a file (one per line)
func readTargetFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var targets []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			targets = append(targets, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return targets, nil
}

// processWebTargetsWithWorkers processes multiple web targets using a worker pool
func processWebTargetsWithWorkers(scanner *web.Scanner, targets []string, languages []string, deep bool, maxDepth int, workers int) ([]*types.ScanResult, error) {
	var allResults []*types.ScanResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Create channels for work distribution
	targetChan := make(chan string, len(targets))
	resultChan := make(chan []*types.ScanResult, len(targets))

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for target := range targetChan {
				log.Info("Scanning target: %s", target)
				results, err := scanner.ScanTarget(target, languages, deep, maxDepth)
				if err != nil {
					log.Warn("Failed to scan target %s: %v", target, err)
					resultChan <- []*types.ScanResult{}
					continue
				}
				resultChan <- results
			}
		}()
	}

	// Send targets to workers
	go func() {
		for _, target := range targets {
			targetChan <- target
		}
		close(targetChan)
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Process results as they come in
	for results := range resultChan {
		mu.Lock()
		allResults = append(allResults, results...)
		mu.Unlock()
	}

	return allResults, nil
}
