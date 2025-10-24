package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	// General settings
	Verbose     bool   `mapstructure:"verbose"`
	Output      string `mapstructure:"output"`
	Format      string `mapstructure:"format"` // json, html, text
	Workers     int    `mapstructure:"workers"`
	Timeout     int    `mapstructure:"timeout"`
	
	// GitHub settings
	GitHubToken string `mapstructure:"github_token"`
	GitHubOrg   string `mapstructure:"github_org"`
	GitHubRepo  string `mapstructure:"github_repo"`
	MaxRepos    int    `mapstructure:"max_repos"`
	
	// Target settings
	Targets     []string `mapstructure:"targets"`
	TargetFile  string   `mapstructure:"target_file"`
	
	// Scanning settings
	SafeSpaces  []string `mapstructure:"safe_spaces"`
	Languages   []string `mapstructure:"languages"`
	DeepScan    bool     `mapstructure:"deep_scan"`
	
	// Rate limiting
	RateLimit   int `mapstructure:"rate_limit"`
	Delay       int `mapstructure:"delay"`
	
	// Output settings
	SaveResults bool   `mapstructure:"save_results"`
	OutputDir   string `mapstructure:"output_dir"`
	
	// Web scanning
	UserAgent   string `mapstructure:"user_agent"`
	FollowRedirects bool `mapstructure:"follow_redirects"`
	
	// Database settings (for AutoAR integration)
	Database struct {
		Type     string `mapstructure:"type"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
	} `mapstructure:"database"`
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		Verbose:         false,
		Output:          "",
		Format:          "text",
		Workers:         10,
		Timeout:         30,
		MaxRepos:        50,
		Languages:       []string{"npm", "pip", "composer", "mvn", "rubygems"},
		DeepScan:        false,
		RateLimit:       100,
		Delay:           100,
		SaveResults:     true,
		OutputDir:       "./results",
		UserAgent:       "Confused-DepConfusion-Scanner/2.0",
		FollowRedirects: true,
		SafeSpaces:      []string{},
	}
}

// Load loads configuration from files and environment variables
func Load(cfg *Config) {
	// Store CLI flag values before they get overwritten
	cliGitHubToken := cfg.GitHubToken
	cliVerbose := cfg.Verbose
	cliOutput := cfg.Output
	cliFormat := cfg.Format
	cliWorkers := cfg.Workers
	cliTimeout := cfg.Timeout
	cliSafeSpaces := cfg.SafeSpaces
	cliOutputDir := cfg.OutputDir
	cliSaveResults := cfg.SaveResults
	cliUserAgent := cfg.UserAgent

	// Set config file
	viper.SetConfigName("confused")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.confused")
	viper.AddConfigPath("/etc/confused")

	// Set environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CONFUSED")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Only warn if it's not a "file not found" error
			// Logger might not be initialized yet, so we can't use it here
		}
	}

	// Unmarshal config
	if err := viper.Unmarshal(cfg); err != nil {
		// Config unmarshal error - this is a critical error
		panic(err)
	}

	// Restore CLI flag values (CLI flags take precedence)
	if cliGitHubToken != "" {
		cfg.GitHubToken = cliGitHubToken
	}
	if cliVerbose {
		cfg.Verbose = cliVerbose
	}
	if cliOutput != "" {
		cfg.Output = cliOutput
	}
	if cliFormat != "" {
		cfg.Format = cliFormat
	}
	if cliWorkers > 0 {
		cfg.Workers = cliWorkers
	}
	if cliTimeout > 0 {
		cfg.Timeout = cliTimeout
	}
	if len(cliSafeSpaces) > 0 {
		cfg.SafeSpaces = cliSafeSpaces
	}
	if cliOutputDir != "" {
		cfg.OutputDir = cliOutputDir
	}
	cfg.SaveResults = cliSaveResults
	if cliUserAgent != "" {
		cfg.UserAgent = cliUserAgent
	}
}

// GetTimeout returns the timeout as a duration
func (c *Config) GetTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

// GetDelay returns the delay as a duration
func (c *Config) GetDelay() time.Duration {
	return time.Duration(c.Delay) * time.Millisecond
}
