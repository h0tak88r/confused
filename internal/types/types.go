package types

import (
	"context"
	"sync"
	"time"
)

// ScanResult represents the result of a dependency confusion scan
type ScanResult struct {
	Target       string                 `json:"target"`
	Type         string                 `json:"type"`
	Language     string                 `json:"language"`
	Vulnerable   []string               `json:"vulnerable_packages"`
	Safe         []string               `json:"safe_packages"`
	Total        int                    `json:"total_packages"`
	Timestamp    time.Time              `json:"timestamp"`
	Duration     time.Duration          `json:"duration"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ScanResults represents multiple scan results
type ScanResults struct {
	Results []ScanResult `json:"results"`
	Summary struct {
		TotalTargets    int `json:"total_targets"`
		VulnerableCount int `json:"vulnerable_count"`
		SafeCount       int `json:"safe_count"`
		Duration        time.Duration `json:"total_duration"`
	} `json:"summary"`
}

// NewScanResult creates a new scan result
func NewScanResult(target, scanType, language string) *ScanResult {
	return &ScanResult{
		Target:    target,
		Type:      scanType,
		Language:  language,
		Vulnerable: []string{},
		Safe:      []string{},
		Metadata:  make(map[string]interface{}),
		Timestamp: time.Now(),
	}
}

// AddVulnerable adds a vulnerable package to the result
func (sr *ScanResult) AddVulnerable(pkg string) {
	sr.Vulnerable = append(sr.Vulnerable, pkg)
}

// AddSafe adds a safe package to the result
func (sr *ScanResult) AddSafe(pkg string) {
	sr.Safe = append(sr.Safe, pkg)
}

// Finalize completes the scan result
func (sr *ScanResult) Finalize() {
	sr.Total = len(sr.Vulnerable) + len(sr.Safe)
	sr.Duration = time.Since(sr.Timestamp)
}

// IsVulnerable returns true if the target has vulnerable packages
func (sr *ScanResult) IsVulnerable() bool {
	return len(sr.Vulnerable) > 0
}

// PackageResolver interface for resolving package information
type PackageResolver interface {
	ReadPackagesFromFile(string) error
	PackagesNotInPublic() []string
	GetPackageCount() int
	GetLanguage() string
}

// EnhancedPackageResolver provides additional functionality for advanced scanning
type EnhancedPackageResolver interface {
	PackageResolver
	SetContext(context.Context)
	SetTimeout(time.Duration)
	SetRateLimit(rate int)
	GetPackageDetails() []PackageDetail
}

// PackageDetail represents detailed information about a package
type PackageDetail struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Type        string            `json:"type"` // dependency, devDependency, etc.
	Vulnerable  bool              `json:"vulnerable"`
	Reason      string            `json:"reason,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// DependencyFile represents a discovered dependency file
type DependencyFile struct {
	URL      string
	Path     string
	Language string
	Content  []byte
	Size     int64
}

// WorkerPool represents a pool of workers for concurrent processing
type WorkerPool struct {
	workers    int
	jobQueue   chan func()
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers:  workers,
		jobQueue: make(chan func(), workers*2),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

// Stop stops the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.jobQueue)
	wp.wg.Wait()
}

// Submit submits a job to the worker pool
func (wp *WorkerPool) Submit(job func()) {
	select {
	case wp.jobQueue <- job:
	case <-wp.ctx.Done():
		return
	}
}

// worker is the worker goroutine
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for {
		select {
		case job, ok := <-wp.jobQueue:
			if !ok {
				return
			}
			job()
		case <-wp.ctx.Done():
			return
		}
	}
}
