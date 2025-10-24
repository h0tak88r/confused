package github

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v58/github"
	"github.com/h0tak88r/confused/internal/resolvers"
	"github.com/h0tak88r/confused/internal/types"
	"github.com/h0tak88r/confused/pkg/config"
	"github.com/h0tak88r/confused/pkg/logger"
	"golang.org/x/oauth2"
)

// Client represents a GitHub API client
type Client struct {
	client *github.Client
	ctx    context.Context
	config *config.Config
	logger *logger.Logger
}

// New creates a new GitHub client
func New(cfg *config.Config, log *logger.Logger) (*Client, error) {
	ctx := context.Background()
	
	var client *github.Client
	
	if cfg.GitHubToken != "" {
		log.Debug("Using GitHub token for authentication")
		log.Info("GitHub token provided - using authenticated requests")
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.GitHubToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		log.Warn("No GitHub token provided, using unauthenticated requests (rate limited)")
		log.Warn("For better performance and higher rate limits, use --github-token flag")
		client = github.NewClient(nil)
	}

	return &Client{
		client: client,
		ctx:    ctx,
		config: cfg,
		logger: log,
	}, nil
}

// ScanRepository scans a specific GitHub repository
func (gc *Client) ScanRepository(repo string, languages []string, safeSpaces []string, deep bool) ([]*types.ScanResult, error) {
	// Parse repository name
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository format: %s (expected owner/repo)", repo)
	}
	
	owner, repoName := parts[0], parts[1]
	
	gc.logger.Info("Scanning repository: %s/%s", owner, repoName)
	
	// Get repository information
	repository, _, err := gc.client.Repositories.Get(gc.ctx, owner, repoName)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}
	
	gc.logger.Debug("Repository: %s", repository.GetFullName())
	gc.logger.Debug("Default branch: %s", repository.GetDefaultBranch())
	
	// Get default branch
	defaultBranch := repository.GetDefaultBranch()
	if defaultBranch == "" {
		defaultBranch = "main"
	}
	
	// Scan for dependency files
	results, err := gc.scanBranch(owner, repoName, defaultBranch, languages, safeSpaces)
	if err != nil {
		return nil, fmt.Errorf("failed to scan default branch: %w", err)
	}
	
	// If deep scan is enabled, scan all branches
	if deep {
		branches, err := gc.getBranches(owner, repoName)
		if err != nil {
			gc.logger.Warn("Failed to get branches for deep scan: %v", err)
		} else {
			for _, branch := range branches {
				if branch == defaultBranch {
					continue // Already scanned
				}
				
				branchResults, err := gc.scanBranch(owner, repoName, branch, languages, safeSpaces)
				if err != nil {
					gc.logger.Warn("Failed to scan branch %s: %v", branch, err)
					continue
				}
				
				results = append(results, branchResults...)
			}
		}
	}
	
	return results, nil
}

// ScanOrganization scans all repositories in a GitHub organization
func (gc *Client) ScanOrganization(org string, languages []string, safeSpaces []string, maxRepos int, deep bool) ([]*types.ScanResult, error) {
	gc.logger.Info("Scanning organization: %s", org)
	
	// Get organization repositories
	repos, err := gc.getOrganizationRepos(org, maxRepos)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization repositories: %w", err)
	}
	
	gc.logger.Info("Found %d repositories", len(repos))
	
	// Create worker pool
	workerPool := types.NewWorkerPool(gc.config.Workers)
	workerPool.Start()
	defer workerPool.Stop()
	
	// Results channel
	resultsChan := make(chan []*types.ScanResult, len(repos))
	var allResults []*types.ScanResult
	
	// Submit jobs
	for _, repo := range repos {
		repo := repo // Capture for closure
		workerPool.Submit(func() {
			repoResults, err := gc.ScanRepository(repo.GetFullName(), languages, safeSpaces, deep)
			if err != nil {
				gc.logger.Warn("Failed to scan repository %s: %v", repo.GetFullName(), err)
				resultsChan <- []*types.ScanResult{}
				return
			}
			resultsChan <- repoResults
		})
	}
	
	// Collect results
	for i := 0; i < len(repos); i++ {
		repoResults := <-resultsChan
		allResults = append(allResults, repoResults...)
	}
	
	return allResults, nil
}

// getOrganizationRepos gets all repositories for an organization
func (gc *Client) getOrganizationRepos(org string, maxRepos int) ([]*github.Repository, error) {
	var allRepos []*github.Repository
	page := 1
	perPage := 100
	
	for {
		repos, resp, err := gc.client.Repositories.ListByOrg(gc.ctx, org, &github.RepositoryListByOrgOptions{
			Type: "public",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
		})
		if err != nil {
			return nil, err
		}
		
		allRepos = append(allRepos, repos...)
		
		// Check if we've reached the limit
		if len(allRepos) >= maxRepos {
			allRepos = allRepos[:maxRepos]
			break
		}
		
		// Check if there are more pages
		if resp.NextPage == 0 {
			break
		}
		
		page = resp.NextPage
	}
	
	return allRepos, nil
}

// getBranches gets all branches for a repository
func (gc *Client) getBranches(owner, repo string) ([]string, error) {
	var branches []string
	page := 1
	perPage := 100
	
	for {
		branchList, resp, err := gc.client.Repositories.ListBranches(gc.ctx, owner, repo, &github.BranchListOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
		})
		if err != nil {
			return nil, err
		}
		
		for _, branch := range branchList {
			branches = append(branches, branch.GetName())
		}
		
		// Check if there are more pages
		if resp.NextPage == 0 {
			break
		}
		
		page = resp.NextPage
	}
	
	return branches, nil
}

// scanBranch scans a specific branch for dependency files
func (gc *Client) scanBranch(owner, repo, branch string, languages []string, safeSpaces []string) ([]*types.ScanResult, error) {
	var results []*types.ScanResult
	
	// Get tree for the branch
	tree, _, err := gc.client.Git.GetTree(gc.ctx, owner, repo, branch, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree for branch %s: %w", branch, err)
	}
	
	// Find dependency files
	dependencyFiles := gc.findDependencyFiles(tree.Entries, languages)
	
	gc.logger.Debug("Found %d dependency files in branch %s", len(dependencyFiles), branch)
	
	// Scan each dependency file
	for _, file := range dependencyFiles {
		result, err := gc.scanDependencyFile(owner, repo, file, safeSpaces)
		if err != nil {
			gc.logger.Warn("Failed to scan dependency file %s: %v", file.GetPath(), err)
			continue
		}
		
		if result != nil {
			results = append(results, result)
		}
	}
	
	return results, nil
}

// findDependencyFiles finds dependency files in the repository tree
func (gc *Client) findDependencyFiles(entries []*github.TreeEntry, languages []string) []*github.TreeEntry {
	var dependencyFiles []*github.TreeEntry
	
	// Define file patterns for each language
	filePatterns := map[string][]string{
		"npm":      {"package.json", "package-lock.json", "yarn.lock"},
		"pip":      {"requirements.txt", "requirements-dev.txt", "setup.py", "pyproject.toml"},
		"composer": {"composer.json", "composer.lock"},
		"mvn":      {"pom.xml"},
		"rubygems": {"Gemfile", "Gemfile.lock", "gems.rb"},
	}
	
	// Collect all patterns for the requested languages
	var allPatterns []string
	for _, lang := range languages {
		if patterns, exists := filePatterns[lang]; exists {
			allPatterns = append(allPatterns, patterns...)
		}
	}
	
	// Find matching files
	for _, entry := range entries {
		if entry.GetType() == "blob" {
			fileName := filepath.Base(entry.GetPath())
			for _, pattern := range allPatterns {
				if fileName == pattern {
					dependencyFiles = append(dependencyFiles, entry)
					break
				}
			}
		}
	}
	
	return dependencyFiles
}

// scanDependencyFile scans a specific dependency file
func (gc *Client) scanDependencyFile(owner, repo string, file *github.TreeEntry, safeSpaces []string) (*types.ScanResult, error) {
	// Get file content
	content, err := gc.getFileContent(owner, repo, file.GetSHA())
	if err != nil {
		return nil, fmt.Errorf("failed to get file content: %w", err)
	}
	
	// Determine language from file extension
	language := gc.getLanguageFromFile(file.GetPath())
	if language == "" {
		return nil, fmt.Errorf("unknown language for file: %s", file.GetPath())
	}
	
	// Create scan result
	result := types.NewScanResult(
		fmt.Sprintf("%s/%s:%s", owner, repo, file.GetPath()),
		"github",
		language,
	)
	
	// Get resolver for the language
	resolver, err := gc.getResolverForLanguage(language)
	if err != nil {
		return nil, fmt.Errorf("failed to get resolver for language %s: %w", language, err)
	}
	
	// Create temporary file
	tempFile, err := gc.createTempFile(content)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile)
	
	// Read packages from file
	if err := resolver.ReadPackagesFromFile(tempFile); err != nil {
		return nil, fmt.Errorf("failed to read packages from file: %w", err)
	}
	
	// Get vulnerable packages
	vulnerablePackages := resolver.PackagesNotInPublic()
	
	// Remove safe spaces
	vulnerablePackages = gc.removeSafe(vulnerablePackages, safeSpaces)
	
	// Add to result
	for _, pkg := range vulnerablePackages {
		result.AddVulnerable(pkg)
	}
	
	// Add metadata
	result.Metadata["file_path"] = file.GetPath()
	result.Metadata["file_sha"] = file.GetSHA()
	result.Metadata["file_size"] = file.GetSize()
	
	// Finalize result
	result.Finalize()
	
	return result, nil
}

// getFileContent gets the content of a file from GitHub
func (gc *Client) getFileContent(owner, repo, sha string) ([]byte, error) {
	blob, _, err := gc.client.Git.GetBlob(gc.ctx, owner, repo, sha)
	if err != nil {
		return nil, err
	}
	
	// Decode base64 content
	content := blob.GetContent()
	
	return []byte(content), nil
}

// getLanguageFromFile determines the language from file path
func (gc *Client) getLanguageFromFile(filePath string) string {
	fileName := filepath.Base(filePath)
	
	// Map file names to languages
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

// createTempFile creates a temporary file with the given content
func (gc *Client) createTempFile(content []byte) (string, error) {
	// Create temporary file
	tempFile, err := os.CreateTemp("", "confused-*.tmp")
	if err != nil {
		return "", err
	}
	
	// Write content
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

// getResolverForLanguage returns a resolver for the given language
func (gc *Client) getResolverForLanguage(language string) (types.PackageResolver, error) {
	return resolvers.GetResolverForLanguageWithVerbose(language, gc.config.Verbose)
}

// removeSafe removes known-safe package names from the slice
func (gc *Client) removeSafe(packages []string, safeSpaces []string) []string {
	if len(safeSpaces) == 0 {
		return packages
	}
	
	retSlice := []string{}
	for _, pkg := range packages {
		ignored := false
		for _, safeSpace := range safeSpaces {
			ok, err := filepath.Match(safeSpace, pkg)
			if err != nil {
				gc.logger.Warn("Encountered an error while trying to match a known-safe namespace %s: %v", safeSpace, err)
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
