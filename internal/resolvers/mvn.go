package resolvers

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/h0tak88r/confused2/internal/types"
)

// MVNLookup represents a collection of maven packages to be tested for dependency confusion.
type MVNLookup struct {
	Packages       []MVNPackage
	Verbose        bool
	ctx            context.Context
	timeout        time.Duration
	rateLimit      int
	packageDetails []types.PackageDetail
}

type MVNPackage struct {
	Group    string
	Artifact string
	Version  string
}

// NewMVNLookup constructs an `MVNLookup` struct and returns it.
func NewMVNLookup(verbose bool) types.PackageResolver {
	return &MVNLookup{
		Packages:       []MVNPackage{},
		Verbose:        verbose,
		ctx:            context.Background(),
		timeout:        30 * time.Second,
		rateLimit:      100,
		packageDetails: []types.PackageDetail{},
	}
}

// ReadPackagesFromFile reads package information from an npm package.json file
//
// Returns any errors encountered
func (n *MVNLookup) ReadPackagesFromFile(filename string) error {
	rawfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if n.Verbose {
		fmt.Print("Checking: filename: " + filename + "\n")
	}

	// Check if file is empty or too small to be a valid POM
	if len(rawfile) < 10 {
		if n.Verbose {
			fmt.Printf("Skipping empty or too small POM file: %s\n", filename)
		}
		return nil
	}

	var project MavenProject
	if err := xml.Unmarshal([]byte(rawfile), &project); err != nil {
		if n.Verbose {
			fmt.Printf("Warning: unable to parse POM file %s: %s\n", filename, err)
		}
		return nil // Return nil instead of crashing
	}

	for _, dep := range project.Dependencies {
		n.Packages = append(n.Packages, MVNPackage{dep.GroupId, dep.ArtifactId, dep.Version})
	}

	for _, dep := range project.Build.Plugins {
		n.Packages = append(n.Packages, MVNPackage{dep.GroupId, dep.ArtifactId, dep.Version})
	}

	for _, build := range project.Profiles {
		for _, dep := range build.Build.Plugins {
			n.Packages = append(n.Packages, MVNPackage{dep.GroupId, dep.ArtifactId, dep.Version})
		}
	}

	return nil
}

// PackagesNotInPublic determines if an npm package does not exist in the public npm package repository.
//
// Returns a slice of strings with any npm packages not in the public npm package repository
func (n *MVNLookup) PackagesNotInPublic() []string {
	notavail := []string{}
	for _, pkg := range n.Packages {
		if !n.isAvailableInPublic(pkg, 0) {
			notavail = append(notavail, pkg.Group+"/"+pkg.Artifact)
		}
	}
	return notavail
}

// isAvailableInPublic determines if an npm package exists in the public npm package repository.
//
// Returns true if the package exists in the public npm package repository.
func (n *MVNLookup) isAvailableInPublic(pkg MVNPackage, retry int) bool {
	if retry > 3 {
		fmt.Printf(" [W] Maximum number of retries exhausted for package: %s\n", pkg.Group)
		return false
	}
	if pkg.Group == "" {
		return true
	}

	group := strings.Replace(pkg.Group, ".", "/", -1)
	if n.Verbose {
		fmt.Print("Checking: https://repo1.maven.org/maven2/" + group + "/ ")
	}
	resp, err := http.Get("https://repo1.maven.org/maven2/" + group + "/")
	if err != nil {
		fmt.Printf(" [W] Error when trying to request https://repo1.maven.org/maven2/"+group+"/ : %s\n", err)
		return false
	}
	defer resp.Body.Close()
	if n.Verbose {
		fmt.Printf("%s\n", resp.Status)
	}
	if resp.StatusCode == http.StatusOK {
		npmResp := NpmResponse{}
		body, _ := ioutil.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &npmResp)
		if npmResp.NotAvailable() {
			if n.Verbose {
				fmt.Printf("[W] Package %s was found, but all its versions are unpublished, making anyone able to takeover the namespace.\n", pkg.Group)
			}
			return false
		}
		return true
	} else if resp.StatusCode == 429 {
		fmt.Printf(" [!] Server responded with 429 (Too many requests), throttling and retrying...\n")
		time.Sleep(10 * time.Second)
		retry = retry + 1
		return n.isAvailableInPublic(pkg, retry)
	}
	return false
}

// GetPackageCount returns the number of packages
func (m *MVNLookup) GetPackageCount() int {
	return len(m.Packages)
}

// GetLanguage returns the language name
func (m *MVNLookup) GetLanguage() string {
	return "mvn"
}

// SetContext sets the context for the resolver
func (m *MVNLookup) SetContext(ctx context.Context) {
	m.ctx = ctx
}

// SetTimeout sets the timeout for requests
func (m *MVNLookup) SetTimeout(timeout time.Duration) {
	m.timeout = timeout
}

// SetRateLimit sets the rate limit for requests
func (m *MVNLookup) SetRateLimit(rate int) {
	m.rateLimit = rate
}

// GetPackageDetails returns detailed information about packages
func (m *MVNLookup) GetPackageDetails() []types.PackageDetail {
	if len(m.packageDetails) == 0 {
		m.buildPackageDetails()
	}
	return m.packageDetails
}

// buildPackageDetails builds detailed package information
func (m *MVNLookup) buildPackageDetails() {
	m.packageDetails = []types.PackageDetail{}
	
	for _, pkg := range m.Packages {
		detail := types.PackageDetail{
			Name:    pkg.Group + ":" + pkg.Artifact,
			Version: pkg.Version,
			Type:    "dependency",
			Metadata: map[string]interface{}{
				"group":    pkg.Group,
				"artifact": pkg.Artifact,
				"version":  pkg.Version,
			},
		}
		
		// Check if package is vulnerable
		if !m.isAvailableInPublic(pkg, 0) {
			detail.Vulnerable = true
			detail.Reason = "Package not available in public Maven repository"
		}
		
		m.packageDetails = append(m.packageDetails, detail)
	}
}