package web

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/h0tak88r/confused/internal/resolvers"
	"github.com/h0tak88r/confused/internal/types"
	"github.com/h0tak88r/confused/pkg/logger"
)

// Scanner represents a web dependency scanner
type Scanner struct {
	client    *http.Client
	logger    *logger.Logger
	userAgent string
	timeout   time.Duration
}

// New creates a new web scanner
func New(log *logger.Logger, userAgent string, timeout time.Duration) *Scanner {
	return &Scanner{
		client: &http.Client{
			Timeout: timeout,
		},
		logger:    log,
		userAgent: userAgent,
		timeout:   timeout,
	}
}

// ScanTarget scans a web target for dependency files
func (s *Scanner) ScanTarget(target string, languages []string, deep bool, maxDepth int) ([]*types.ScanResult, error) {
	var results []*types.ScanResult
	
	// Normalize target URL
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "https://" + target
	}
	
	// Parse URL
	u, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	
	s.logger.Info("Scanning web target: %s", target)
	
	// Define common dependency file paths
	filePaths := s.getDependencyFilePaths(languages)
	
	// Scan for dependency files
	for _, filePath := range filePaths {
		result, err := s.scanDependencyFile(u, filePath, languages)
		if err != nil {
			s.logger.Debug("Failed to scan %s: %v", filePath, err)
			continue
		}
		
		if result != nil {
			results = append(results, result)
		}
	}
	
	// If deep scan is enabled, try additional discovery methods
	if deep {
		deepResults, err := s.deepScan(u, languages, maxDepth)
		if err != nil {
			s.logger.Warn("Deep scan failed: %v", err)
		} else {
			results = append(results, deepResults...)
		}
	}
	
	return results, nil
}

// getDependencyFilePaths returns common paths for dependency files
func (s *Scanner) getDependencyFilePaths(languages []string) []string {
	var paths []string
	
	// Common root paths
	commonPaths := []string{
		"package.json",
		"package-lock.json",
		"yarn.lock",
		"requirements.txt",
		"requirements-dev.txt",
		"setup.py",
		"pyproject.toml",
		"composer.json",
		"composer.lock",
		"pom.xml",
		"Gemfile",
		"Gemfile.lock",
		"gems.rb",
	}
	
	// Add language-specific paths
	languagePaths := map[string][]string{
		"npm":      {"package.json", "package-lock.json", "yarn.lock"},
		"pip":      {"requirements.txt", "requirements-dev.txt", "setup.py", "pyproject.toml"},
		"composer": {"composer.json", "composer.lock"},
		"mvn":      {"pom.xml"},
		"rubygems": {"Gemfile", "Gemfile.lock", "gems.rb"},
	}
	
	// If specific languages are requested, only include those paths
	if len(languages) > 0 {
		paths = []string{}
		for _, lang := range languages {
			if langPaths, exists := languagePaths[lang]; exists {
				paths = append(paths, langPaths...)
			}
		}
	} else {
		paths = commonPaths
	}
	
	return paths
}

// scanDependencyFile scans a specific dependency file
func (s *Scanner) scanDependencyFile(baseURL *url.URL, filePath string, languages []string) (*types.ScanResult, error) {
	// Construct full URL
	fileURL := *baseURL
	fileURL.Path = filepath.Join(baseURL.Path, filePath)
	
	s.logger.Debug("Checking: %s", fileURL.String())
	
	// Make HTTP request
	req, err := http.NewRequest("GET", fileURL.String(), nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("User-Agent", s.userAgent)
	
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// Check if file exists and is accessible
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("file not found (status: %d)", resp.StatusCode)
	}
	
	// Read file content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	// Determine language from file path
	language := s.getLanguageFromFile(filePath)
	if language == "" {
		return nil, fmt.Errorf("unknown language for file: %s", filePath)
	}
	
	// Create scan result
	result := types.NewScanResult(
		fmt.Sprintf("%s:%s", baseURL.Host, filePath),
		"web",
		language,
	)
	
	// Get resolver for the language
	resolver, err := s.getResolverForLanguage(language)
	if err != nil {
		return nil, fmt.Errorf("failed to get resolver for language %s: %w", language, err)
	}
	
	// Create temporary file
	tempFile, err := s.createTempFile(content)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer s.cleanupTempFile(tempFile)
	
	// Read packages from file
	if err := resolver.ReadPackagesFromFile(tempFile); err != nil {
		s.logger.Debug("Failed to parse %s: %v", filePath, err)
		return nil, err
	}
	
	// Get vulnerable packages
	vulnerablePackages := resolver.PackagesNotInPublic()
	
	// Add to result
	for _, pkg := range vulnerablePackages {
		result.AddVulnerable(pkg)
	}
	
	// Add metadata
	result.Metadata["file_path"] = filePath
	result.Metadata["file_url"] = fileURL.String()
	result.Metadata["file_size"] = len(content)
	result.Metadata["status_code"] = resp.StatusCode
	
	// Finalize result
	result.Finalize()
	
	return result, nil
}

// deepScan performs additional discovery methods
func (s *Scanner) deepScan(baseURL *url.URL, languages []string, maxDepth int) ([]*types.ScanResult, error) {
	var results []*types.ScanResult
	
	// Try common directory paths
	commonDirs := []string{
		"src/", "lib/", "app/", "web/", "public/", "static/",
		"api/", "backend/", "frontend/", "client/", "server/",
	}
	
	for _, dir := range commonDirs {
		dirResults, err := s.scanDirectory(baseURL, dir, languages)
		if err != nil {
			s.logger.Debug("Failed to scan directory %s: %v", dir, err)
			continue
		}
		results = append(results, dirResults...)
	}
	
	return results, nil
}

// scanDirectory scans a directory for dependency files
func (s *Scanner) scanDirectory(baseURL *url.URL, dir string, languages []string) ([]*types.ScanResult, error) {
	var results []*types.ScanResult
	
	filePaths := s.getDependencyFilePaths(languages)
	
	for _, filePath := range filePaths {
		fullPath := filepath.Join(dir, filePath)
		result, err := s.scanDependencyFile(baseURL, fullPath, languages)
		if err != nil {
			continue // File not found in this directory
		}
		
		if result != nil {
			results = append(results, result)
		}
	}
	
	return results, nil
}

// getLanguageFromFile determines the language from file path
func (s *Scanner) getLanguageFromFile(filePath string) string {
	fileName := filepath.Base(filePath)
	
	fileLanguageMap := map[string]string{
		"package.json":     "npm",
		"package-lock.json": "npm",
		"yarn.lock":        "npm",
		"requirements.txt": "pip",
		"requirements-dev.txt": "pip",
		"setup.py":         "pip",
		"pyproject.toml":   "pip",
		"composer.json":    "composer",
		"composer.lock":    "composer",
		"pom.xml":          "mvn",
		"Gemfile":          "rubygems",
		"Gemfile.lock":     "rubygems",
		"gems.rb":          "rubygems",
	}
	
	return fileLanguageMap[fileName]
}

// getResolverForLanguage returns the appropriate resolver for a language
func (s *Scanner) getResolverForLanguage(language string) (types.PackageResolver, error) {
	return resolvers.GetResolverForLanguage(language)
}

// createTempFile creates a temporary file with the given content
func (s *Scanner) createTempFile(content []byte) (string, error) {
	// Create temporary file
	tempFile, err := os.CreateTemp("", "confused-web-*.tmp")
	if err != nil {
		return "", err
	}
	
	// Write content to file
	if _, err := tempFile.Write(content); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return "", err
	}
	
	// Close file
	if err := tempFile.Close(); err != nil {
		os.Remove(tempFile.Name())
		return "", err
	}
	
	return tempFile.Name(), nil
}

// cleanupTempFile removes a temporary file
func (s *Scanner) cleanupTempFile(filename string) {
	if filename != "" {
		os.Remove(filename)
	}
}
