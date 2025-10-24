# Confused - Advanced Dependency Confusion Scanner

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-2.2.0-orange.svg)](https://github.com/h0tak88r/confused/releases)

An advanced, professional-grade dependency confusion scanner that can scan local files, GitHub repositories, organizations, and web targets for dependency confusion vulnerabilities. Built with Go for high performance and concurrent processing.

## 🚀 Features

### Core Functionality
- **Multi-Language Support**: npm, pip, composer, mvn, rubygems
- **Concurrent Processing**: Configurable worker pools with `-w` flag
- **Multiple Target Types**: Local files, GitHub repos/orgs, web targets
- **Advanced Discovery**: Brute force file discovery for web targets
- **Professional Reporting**: JSON, HTML, and text output formats

### GitHub Integration
- **Repository Scanning**: Scan individual GitHub repositories
- **Organization Scanning**: Scan entire GitHub organizations
- **Deep Scanning**: Scan all branches for comprehensive coverage
- **Rate Limiting**: Intelligent API rate limiting and retry mechanisms

### Web Scanning
- **File Discovery**: Brute force discovery of dependency files
- **Custom Wordlists**: Support for custom discovery wordlists
- **Deep Scanning**: Configurable depth for comprehensive coverage
- **Smart Detection**: Validates file content before processing

### Professional Features
- **Structured Logging**: Color-coded, configurable logging levels
- **Configuration Management**: YAML configuration files
- **Database Integration**: Ready for AutoAR framework integration
- **Comprehensive Reporting**: Detailed vulnerability reports with metadata

## 📋 What is Dependency Confusion?

On February 9th, 2021, security researcher Alex Birsan [published groundbreaking research](https://medium.com/@alex.birsan/dependency-confusion-4a5d60fec610) on dependency confusion attacks. This attack vector exploits the resolution order of package managers, allowing attackers to inject malicious packages into applications that use private package repositories.

Microsoft [released a comprehensive whitepaper](https://azure.microsoft.com/en-gb/resources/3-ways-to-mitigate-risk-using-private-package-feeds/) describing mitigation strategies, though the root cause remains in many ecosystems.

## 🛠️ Installation

### Pre-built Binaries
Download the latest release from the [releases page](https://github.com/h0tak88r/confused/releases/latest).

### From Source
```bash
git clone https://github.com/h0tak88r/confused
cd confused
go mod tidy
go build -o confused
```

### Go Install
```bash
go install github.com/h0tak88r/confused@latest
```

## 🚀 Quick Start

### Scan a Local File
```bash
# Scan a package.json file
confused scan package.json

# Scan with specific language and workers
confused scan -l pip requirements.txt -w 20

# Scan with safe spaces (known secure namespaces)
confused scan -l npm package.json --safe-spaces "@mycompany/*,@trusted/*"
```

### Scan GitHub Repository
```bash
# Scan a specific repository
confused github repo microsoft/PowerShell

# Scan with custom settings
confused github repo facebook/react -w 20 --languages npm,pip --deep
```

### Scan GitHub Organization
```bash
# Scan an organization (default: 50 repos)
confused github org microsoft

# Scan with custom limits
confused github org google --max-repos 100 -w 30
```

### Scan Web Target
```bash
# Scan a single web application
confused web https://example.com

# Scan multiple targets with high concurrency
confused web https://site1.com https://site2.com https://site3.com --workers 20

# Scan from target file
confused web --target-file targets.txt --workers 15 --deep

# Deep scan with custom wordlist
confused web example.com --deep --wordlist custom-wordlist.txt
```

## 📖 Detailed Usage

### Command Structure
```
confused [command] [flags] [arguments]

Commands:
  scan      Scan a local dependency file
  github    GitHub scanning commands
  web       Web target scanning
  config    Configuration management

Flags:
  -v, --verbose          Verbose output
  -o, --output string    Output file path
  -f, --format string    Output format (text, json, html)
  -w, --workers int      Number of concurrent workers (default 10)
  --timeout int          Request timeout in seconds (default 30)
  --safe-spaces strings  Known-safe namespaces (supports wildcards)
  --output-dir string    Output directory for results (default "./results")
  --save                 Save results to files (default true)
  --github-token string  GitHub API token
  --user-agent string    User agent for HTTP requests
```

### Scan Command
```bash
# Basic usage
confused scan <file>

# Options
  -l, --language string     Package manager (npm, pip, composer, mvn, rubygems)
  --safe-spaces strings     Known-safe namespaces (supports wildcards)
```

### GitHub Commands
```bash
# Repository scanning
confused github repo <owner/repo> [flags]
  --languages strings       Package managers to scan for
  --safe-spaces strings     Known-safe namespaces
  --deep                    Perform deep scan including all branches
  -w, --workers int         Number of concurrent workers (default 10)

# Organization scanning
confused github org <organization> [flags]
  --languages strings       Package managers to scan for
  --safe-spaces strings     Known-safe namespaces
  --max-repos int          Maximum number of repositories to scan
  --deep                    Perform deep scan including all branches
  -w, --workers int         Number of concurrent workers (default 10)
```

### Web Command
```bash
# Single target
confused web <target> [flags]

# Multiple targets
confused web <target1> <target2> <target3> [flags]

# Target file
confused web --target-file targets.txt [flags]

Options:
  --languages strings       Package managers to scan for
  --safe-spaces strings     Known-safe namespaces
  --deep                    Perform deep scan with extensive file discovery
  --wordlist strings        Custom wordlist for file discovery
  --max-depth int          Maximum directory depth for discovery
  --target-file string      File containing list of targets (one per line)
  -w, --workers int         Number of concurrent workers (default 10)
```

## ⚙️ Configuration

### Configuration File
Create a `confused.yaml` file in your project root, home directory, or `/etc/confused/`:

```yaml
# General settings
verbose: false
output: ""
format: "text"  # text, json, html
workers: 10
timeout: 30

# GitHub settings
github_token: ""  # Set via CONFUSED_GITHUB_TOKEN environment variable
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
```

### Environment Variables
Set environment variables with the `CONFUSED_` prefix:
```bash
export CONFUSED_GITHUB_TOKEN="your_token_here"
export CONFUSED_WORKERS="20"
export CONFUSED_VERBOSE="true"
```

### Generate Sample Configuration
```bash
confused config generate -o confused.yaml
```

## 🚀 Advanced Features

### Concurrent Processing
All commands support the `-w` (workers) flag for concurrent processing:

```bash
# High-performance web scanning
confused web --target-file targets.txt --workers 50 --verbose

# Fast GitHub organization scanning
confused github org microsoft --workers 25 --max-repos 200

# Concurrent repository scanning
confused github repo facebook/react --workers 15 --deep
```

### Multiple Target Support
Web scanning supports multiple targets in various ways:

```bash
# Command line targets
confused web https://site1.com https://site2.com https://site3.com

# Target file (one per line, comments supported)
confused web --target-file targets.txt

# Mixed approach
confused web https://site1.com --target-file additional-targets.txt
```

### Target File Format
Create a `targets.txt` file with your targets:
```bash
# This is a comment - will be ignored
https://example.com
https://app.example.com
https://api.example.com
# Another comment
https://admin.example.com
```

### Performance Optimization
- **Workers**: Adjust based on your system (CPU cores, memory, network)
- **Rate Limiting**: Built-in GitHub API rate limiting
- **Concurrent Processing**: Multiple targets processed simultaneously
- **Smart Discovery**: Efficient file discovery algorithms

## 📊 Output Formats

### Text Output (Default)
```
[INFO] Starting dependency confusion scan...
[INFO] Target: package.json
[INFO] Language: npm
[INFO] Workers: 10
[WARN] Issues found, the following packages are not available in public package repositories:
[WARN]  [!] internal-package-1
[WARN]  [!] @company/private-package
```

### JSON Output
```bash
confused scan package.json -f json -o results.json
```

```json
{
  "target": "package.json",
  "type": "file",
  "language": "npm",
  "vulnerable_packages": ["internal-package-1", "@company/private-package"],
  "safe_packages": ["lodash", "express"],
  "total_packages": 4,
  "timestamp": "2024-01-01T12:00:00Z",
  "duration": "2.5s",
  "metadata": {
    "file_path": "package.json",
    "file_size": 1024
  }
}
```

### HTML Output
```bash
confused scan package.json -f html -o report.html
```

## 🔧 Advanced Usage

### Concurrent Processing
```bash
# Use 50 workers for high-performance scanning
confused github org microsoft -w 50

# Scan multiple targets in parallel
confused web https://example.com -w 20
confused web https://app.example.com -w 20
```

### Custom Wordlists
```bash
# Create custom wordlist
echo -e "admin\napi\napp\nbackup" > custom-wordlist.txt

# Use custom wordlist for web scanning
confused web https://target.com --wordlist custom-wordlist.txt --deep
```

### Safe Spaces Configuration
```bash
# Exclude known safe namespaces
confused scan package.json --safe-spaces "@mycompany/*,@trusted/*,internal-*"

# Use wildcards for pattern matching
confused scan composer.json --safe-spaces "company-*"
```

## 🔍 GitHub Integration

### Repository Scanning
```bash
# Scan a specific repository
confused github repo microsoft/PowerShell

# Deep scan with all branches
confused github repo facebook/react --deep

# Scan with specific languages
confused github repo google/go --languages npm,pip
```

### Organization Scanning
```bash
# Scan entire organization
confused github org microsoft

# Limit number of repositories
confused github org google --max-repos 100

# Deep scan with high concurrency
confused github org facebook --deep -w 30
```

### GitHub Token Setup
```bash
# Set GitHub token for enhanced API access
export CONFUSED_GITHUB_TOKEN="ghp_your_token_here"

# Or use command line flag
confused github org microsoft --github-token "ghp_your_token_here"
```

## 🌐 Web Scanning

### Basic Web Scanning
```bash
# Scan a web application
confused web https://example.com

# Scan with specific languages
confused web https://app.example.com --languages npm,pip
```

### Deep Web Scanning
```bash
# Deep scan with extensive discovery
confused web https://target.com --deep

# Custom depth and wordlist
confused web https://target.com --deep --max-depth 5 --wordlist custom.txt
```

### File Discovery
The web scanner automatically discovers dependency files through:
- Common file paths (`/package.json`, `/requirements.txt`, etc.)
- Directory brute forcing
- Sitemap analysis
- Response content analysis

## 📈 Performance Optimization

### Worker Configuration
```bash
# High-performance scanning
confused github org microsoft -w 50

# Memory-constrained environment
confused scan package.json -w 5
```

### Rate Limiting
```bash
# Conservative rate limiting
confused web https://target.com --delay 500

# Aggressive scanning (use responsibly)
confused web https://target.com --delay 50
```

## 🔒 Security Considerations

### Safe Spaces
Always configure safe spaces for known secure namespaces:
```bash
confused scan package.json --safe-spaces "@mycompany/*,@trusted/*"
```

### Rate Limiting
Respect API rate limits to avoid service disruption:
```bash
# Conservative approach
confused github org microsoft --delay 200

# Check API limits before scanning
confused config validate
```

### Legal Compliance
- Ensure you have proper authorization before scanning targets
- Respect robots.txt and rate limits
- Use responsibly and ethically

## 🤝 Integration with AutoAR

This tool is designed to integrate seamlessly with the [AutoAR framework](https://github.com/h0tak88r/AutoAR):

```bash
# Use as AutoAR module
/app/main.sh confused scan -d example.com
/app/main.sh confused github org -o microsoft
/app/main.sh confused web -d example.com --deep
```

### Database Integration
Configure database settings for AutoAR integration:
```yaml
database:
  type: "postgresql"
  host: "localhost"
  port: 5432
  username: "confused"
  password: "password"
  database: "autoar"
```

## 🐛 Troubleshooting

### Common Issues

**GitHub API Rate Limiting**
```bash
# Use GitHub token for higher rate limits
export CONFUSED_GITHUB_TOKEN="your_token"

# Reduce concurrency
confused github org microsoft -w 5
```

**Web Scanning Timeouts**
```bash
# Increase timeout
confused web https://target.com --timeout 60

# Reduce workers
confused web https://target.com -w 5
```

**Memory Issues**
```bash
# Reduce workers and concurrency
confused github org microsoft -w 5
```

### Debug Mode
```bash
# Enable verbose logging
confused scan package.json -v

# Check configuration
confused config validate
```

## 📝 Examples

### Complete Workflow
```bash
# 1. Generate configuration
confused config generate -o confused.yaml

# 2. Configure GitHub token
export CONFUSED_GITHUB_TOKEN="your_token"

# 3. Scan local file
confused scan package.json -v

# 4. Scan GitHub organization
confused github org microsoft --max-repos 20 -w 15

# 5. Scan web target
confused web https://example.com --deep -w 10

# 6. Generate HTML report
confused scan package.json -f html -o report.html
```

### Batch Processing
```bash
# Process multiple files
for file in *.json; do
  confused scan "$file" -f json -o "results/${file%.json}.json"
done

# Process multiple targets
while read target; do
  confused web "$target" --deep -w 5
done < targets.txt
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Alex Birsan** for the original dependency confusion research
- **Microsoft** for the comprehensive mitigation whitepaper
- **ProjectDiscovery** for inspiration and tooling
- **The Go Community** for excellent libraries and tools

## 📋 Changelog

### v2.2.0 - Simplified GitHub Actions & Clean Project
**Release Date**: October 24, 2025

#### 🚀 New Features
- **Simplified GitHub Actions**: Streamlined CI/CD with simple release workflow
- **Automatic Releases**: GitHub Action automatically creates releases on version tags
- **Clean Project Structure**: Removed complex workflows and unnecessary files

#### 🔧 Improvements
- **GitHub Actions**: Replaced complex CodeQL and linting workflows with simple release workflow
- **Build Process**: Simplified to single Linux binary builds
- **Documentation**: Updated project overview and structure documentation

#### 🗑️ Removed
- **CodeQL Analysis**: Removed complex security scanning workflow
- **GoLint Workflow**: Removed golangci-lint workflow
- **GoReleaser**: Removed complex cross-platform build configuration
- **Test Artifacts**: Cleaned up all testing and temporary files

### v2.1.0 - Enhanced Concurrency & Multi-Target Support
**Release Date**: October 24, 2025

#### 🚀 New Features
- **Multiple Target Support**: Web scanning now supports multiple targets via command line or target file
- **Target File Support**: `--target-file` flag for reading targets from file (supports comments)
- **Universal Workers Flag**: All commands now support `-w/--workers` flag for concurrent processing
- **Enhanced Web Scanner**: Complete rewrite with proper HTTP handling and file discovery
- **Concurrent Processing**: Worker pool pattern for efficient parallel target processing

#### 🔧 Improvements
- **Better Error Handling**: Graceful handling of failed targets without stopping entire scan
- **Performance Optimization**: Significant speed improvements with concurrent processing
- **Flexible Arguments**: Smart argument validation for target file vs command line targets
- **Enhanced Logging**: Better progress reporting for concurrent operations

#### 🐛 Bug Fixes
- **Maven Parser**: Fixed crashes on malformed `pom.xml` files
- **Rate Limiting**: Improved GitHub API rate limit handling
- **File Path Issues**: Fixed URL normalization in web scanning
- **Memory Management**: Better temporary file cleanup

### v2.0.0 - Professional Go Project Structure
**Release Date**: October 24, 2025

#### 🏗️ Major Restructuring
- **Standard Go Layout**: Complete restructure following Go best practices
- **Cobra CLI**: Professional command-line interface with subcommands
- **Viper Configuration**: Flexible configuration management
- **Structured Logging**: Color-coded, configurable logging system
- **Package Architecture**: Clean separation of concerns with `cmd/`, `pkg/`, `internal/`

#### 🚀 New Features
- **GitHub Integration**: Full GitHub API integration with OAuth2
- **Web Scanning**: HTTP-based target scanning with file discovery
- **Multiple Languages**: Support for npm, pip, composer, mvn, rubygems
- **Advanced Reporting**: JSON, HTML, and text output formats
- **Configuration Management**: YAML configuration files
- **Docker Support**: Containerized deployment ready

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/h0tak88r/confused/issues)
- **Discussions**: [GitHub Discussions](https://github.com/h0tak88r/confused/discussions)
- **Documentation**: [Wiki](https://github.com/h0tak88r/confused/wiki)

---

**⚠️ Disclaimer**: This tool is for educational and authorized testing purposes only. Always ensure you have proper authorization before scanning any target. The authors are not responsible for any misuse of this tool.