# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.2.0] - 2025-10-24

### Added
- **Simplified GitHub Actions**: Streamlined CI/CD with simple release workflow
- **Automatic Releases**: GitHub Action automatically creates releases on version tags
- **Clean Project Structure**: Removed complex workflows and unnecessary files

### Changed
- **GitHub Actions**: Replaced complex CodeQL and linting workflows with simple release workflow
- **Build Process**: Simplified to single Linux binary builds
- **Documentation**: Updated project overview and structure documentation

### Removed
- **CodeQL Analysis**: Removed complex security scanning workflow
- **GoLint Workflow**: Removed golangci-lint workflow
- **GoReleaser**: Removed complex cross-platform build configuration
- **Test Artifacts**: Cleaned up all testing and temporary files

### Technical Details
- **Release Workflow**: Simple workflow that triggers on version tags (v*)
- **Build Process**: Uses Go build with version and date injection
- **Release Creation**: Uses softprops/action-gh-release for automatic releases
- **Project Cleanup**: Removed 29 files including test artifacts and documentation

## [2.1.0] - 2025-10-24

### Added
- **Multiple Target Support**: Web scanning now supports multiple targets via command line arguments
- **Target File Support**: `--target-file` flag for reading targets from file (one per line, comments supported)
- **Universal Workers Flag**: All commands now support `-w/--workers` flag for concurrent processing
- **Enhanced Web Scanner**: Complete rewrite with proper HTTP handling and file discovery
- **Concurrent Processing**: Worker pool pattern for efficient parallel target processing
- **Smart Argument Validation**: Flexible argument handling for target file vs command line targets

### Changed
- **Web Command**: Now accepts multiple targets or target file instead of single target only
- **Performance**: Significant speed improvements with concurrent processing across all commands
- **Error Handling**: Graceful handling of failed targets without stopping entire scan
- **Logging**: Enhanced progress reporting for concurrent operations

### Fixed
- **Maven Parser**: Fixed crashes on malformed `pom.xml` files (replaced `log.Fatalf` with graceful error handling)
- **Rate Limiting**: Improved GitHub API rate limit handling and error messages
- **File Path Issues**: Fixed URL normalization in web scanning (proper `https://` handling)
- **Memory Management**: Better temporary file cleanup in web scanner
- **Argument Validation**: Fixed web command argument validation for target file usage

### Technical Details
- **Web Scanner**: Complete rewrite with proper HTTP client, temporary file handling, and concurrent processing
- **Worker Pool**: Implemented efficient worker pool pattern for all scanning commands
- **File Discovery**: Enhanced file discovery with proper error handling and content validation
- **Concurrency**: Added `sync.WaitGroup` and channels for safe concurrent processing

## [2.0.0] - 2024-01-01

### Added
- **Professional CLI Framework**: Complete rewrite using Cobra CLI framework
- **GitHub Integration**: 
  - Repository scanning with `confused github repo <owner/repo>`
  - Organization scanning with `confused github org <organization>`
  - Deep scanning across all branches
  - Rate limiting and retry mechanisms
- **Web Scanning**: 
  - Web target scanning with `confused web <target>`
  - Brute force file discovery
  - Custom wordlist support
  - Smart file content validation
- **Concurrent Processing**: 
  - Worker pool pattern with `-w` flag
  - Configurable concurrency levels
  - Parallel target processing
- **Advanced Configuration**:
  - YAML configuration files
  - Environment variable support
  - Configuration validation
- **Professional Logging**:
  - Color-coded structured logging
  - Multiple log levels (DEBUG, INFO, WARN, ERROR, FATAL)
  - File logging support
- **Comprehensive Reporting**:
  - JSON output format
  - HTML report generation
  - Detailed vulnerability metadata
  - Batch processing support
- **Enhanced Package Resolvers**:
  - Extended interfaces with detailed package information
  - Context and timeout support
  - Rate limiting per resolver
- **AutoAR Integration**:
  - Database integration ready
  - Framework-compatible APIs
  - Docker containerization
- **Docker Support**:
  - Multi-stage Dockerfile
  - Docker Compose configuration
  - Health checks and security hardening
- **Build System**:
  - Makefile for easy building
  - Cross-compilation support
  - Automated dependency management
- **Documentation**:
  - Comprehensive README with examples
  - API documentation
  - Integration guides

### Changed
- **Architecture**: Complete rewrite from simple CLI to professional framework
- **Package Structure**: Modular design with separate packages for different functionalities
- **Error Handling**: Improved error handling and user feedback
- **Performance**: Significant performance improvements with concurrent processing
- **User Experience**: Professional CLI with help system and validation

### Deprecated
- Legacy command-line interface (replaced by new CLI framework)
- Simple file-only scanning (now supports multiple target types)

### Removed
- Basic flag-based argument parsing
- Single-threaded processing
- Limited output formats

### Fixed
- Memory leaks in long-running scans
- Race conditions in concurrent operations
- Improved error messages and debugging

### Security
- Non-root container execution
- Input validation and sanitization
- Rate limiting to prevent abuse
- Safe file handling and cleanup

## [1.0.0] - 2021-02-09

### Added
- Initial release
- Basic dependency confusion detection
- Support for npm, pip, composer, mvn, rubygems
- Simple CLI interface
- Basic file scanning functionality

### Changed
- N/A

### Deprecated
- N/A

### Removed
- N/A

### Fixed
- N/A

### Security
- Basic input validation

## [0.4] - 2021-02-09

### Added
- npm: In case package was found, also check if all the package versions have been unpublished. This makes the package vulnerable to takeover
- npm: Check for http & https and GitHub version references
- MVN (Maven) support

### Changed
- Fixed a bug where the pip requirements.txt parser processes a 'tilde equals' sign.
- Fixed an issue that would detect git repository urls as matches

## [0.3] - 2021-02-09

### Added
- PHP (composer) support
- Command line parameter to let the user to flag namespaces as known-safe

### Changed
- Python (pypi) dependency definition files that use line continuation are now parsed correctly
- Revised the output to clarify the usage
- Fixed npm package.json file parsing issues when the source file is not following the specification

## [0.2] - 2021-02-09

### Changed
- npm registry checkup url
- Throttle the rate of requests in case of 429 (Too many requests) responses

## [0.1] - 2021-02-09

### Added
- Initial release