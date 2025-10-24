# Confused - Project Overview

## 📁 Project Structure

```
confused/
├── .github/                    # GitHub Actions workflows
│   └── workflows/
│       └── release.yml         # Simple release workflow
├── cmd/                        # Main application
│   └── confused/
│       ├── main.go            # Application entry point
│       └── commands.go        # CLI commands and handlers
├── internal/                   # Private application code
│   ├── resolvers/             # Package manager resolvers
│   │   ├── composer.go        # Composer resolver
│   │   ├── factory.go         # Resolver factory
│   │   ├── mvn.go             # Maven resolver
│   │   ├── mvnparser.go       # Maven XML parser
│   │   ├── npm.go             # NPM resolver
│   │   ├── pip.go             # Pip resolver
│   │   ├── rubygems.go        # RubyGems resolver
│   │   └── util.go            # Utility functions
│   └── types/                 # Type definitions
│       └── types.go           # Core types and interfaces
├── pkg/                       # Public library code
│   ├── config/                # Configuration management
│   │   └── config.go
│   ├── github/                # GitHub API integration
│   │   └── github.go
│   ├── logger/                # Logging system
│   │   └── logger.go
│   └── web/                   # Web scanning
│       └── web.go
├── .gitignore                 # Git ignore rules
├── CHANGELOG.md              # Project changelog
├── Dockerfile                # Docker configuration
├── docker-compose.yml        # Docker Compose setup
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── LICENSE                   # MIT License
├── Makefile                  # Build automation
├── README.md                 # Project documentation
└── confused.yaml             # Sample configuration
```

## 🚀 Key Features

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

## 🛠️ Build & Development

### Prerequisites
- Go 1.21+
- Git

### Building
```bash
# Build binary
make build

# Cross-compile for multiple platforms
make cross-compile

# Clean build artifacts
make clean

# Install to GOPATH/bin
make install
```

### Docker
```bash
# Build Docker image
docker build -t confused .

# Run with Docker Compose
docker-compose up
```

## 📦 Distribution

### GitHub Actions
Simple automated releases with:
- Build on version tags (v*)
- Automatic GitHub release creation
- Single Linux binary distribution

## 🔧 Configuration

### Configuration File
Create `confused.yaml` in project root, home directory, or `/etc/confused/`:

```yaml
# General settings
verbose: false
output: ""
format: "text"  # text, json, html
workers: 10
timeout: 30

# GitHub settings
github_token: ""  # Set via CONFUSED_GITHUB_TOKEN
max_repos: 50

# Scanning settings
safe_spaces: []  # Known-safe namespaces
languages: ["npm", "pip", "composer", "mvn", "rubygems"]
deep_scan: false

# Output settings
save_results: true
output_dir: "./results"
```

## 📊 Usage Examples

### Local File Scanning
```bash
./confused scan package.json
./confused scan requirements.txt -l pip
```

### GitHub Scanning
```bash
./confused github repo facebook/react --workers 15
./confused github org microsoft --workers 25 --max-repos 100
```

### Web Scanning
```bash
./confused web https://example.com
./confused web https://site1.com https://site2.com --workers 20
./confused web --target-file targets.txt --deep
```

## 🎯 Ready for GitHub

The project is now clean and ready for GitHub publication with:
- ✅ Professional Go project structure
- ✅ Comprehensive documentation
- ✅ CI/CD workflows
- ✅ Docker support
- ✅ Automated releases
- ✅ Clean codebase (no test artifacts)
- ✅ Proper .gitignore
- ✅ MIT License
