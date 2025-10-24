package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/h0tak88r/confused2/internal/types"
)

// ComposerJSON represents the dependencies of a composer package
type ComposerJSON struct {
	Require    map[string]string `json:"require,omitempty"`
	RequireDev map[string]string `json:"require-dev,omitempty"`
}

// ComposerLookup represents a collection of composer packages to be tested for dependency confusion.
type ComposerLookup struct {
	Packages       []ComposerPackage
	Verbose        bool
	ctx            context.Context
	timeout        time.Duration
	rateLimit      int
	packageDetails []types.PackageDetail
}

type ComposerPackage struct {
	Name    string
	Version string
}

// NewComposerLookup constructs a `ComposerLookup` struct and returns it.
func NewComposerLookup(verbose bool) types.PackageResolver {
	return &ComposerLookup{
		Packages:       []ComposerPackage{},
		Verbose:        verbose,
		ctx:            context.Background(),
		timeout:        30 * time.Second,
		rateLimit:      100,
		packageDetails: []types.PackageDetail{},
	}
}

// ReadPackagesFromFile reads package information from a composer `composer.json` file
//
// Returns any errors encountered
func (c *ComposerLookup) ReadPackagesFromFile(filename string) error {
	rawfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	data := ComposerJSON{}
	err = json.Unmarshal([]byte(rawfile), &data)
	if err != nil {
		fmt.Printf(" [W] Non-fatal issue encountered while reading %s : %s\n", filename, err)
	}
	for pkgname, pkgversion := range data.Require {
		c.Packages = append(c.Packages, ComposerPackage{pkgname, pkgversion})
	}
	for pkgname, pkgversion := range data.RequireDev {
		c.Packages = append(c.Packages, ComposerPackage{pkgname, pkgversion})
	}
	return nil
}

// PackagesNotInPublic determines if a composer package does not exist in the public composer package repository.
//
// Returns a slice of strings with any composer packages not in the public composer package repository
func (c *ComposerLookup) PackagesNotInPublic() []string {
	notavail := []string{}
	for _, pkg := range c.Packages {
		if c.localReference(pkg.Version) || c.urlReference(pkg.Version) || c.gitReference(pkg.Version) {
			continue
		}
		if !c.isAvailableInPublic(pkg.Name, 0) {
			notavail = append(notavail, pkg.Name)
		}
	}
	return notavail
}

// isAvailableInPublic determines if a composer package exists in the public composer package repository.
//
// Returns true if the package exists in the public composer package repository.
func (c *ComposerLookup) isAvailableInPublic(pkgname string, retry int) bool {
	if retry > 3 {
		fmt.Printf(" [W] Maximum number of retries exhausted for package: %s\n", pkgname)
		return false
	}
	if c.Verbose {
		fmt.Print("Checking: https://packagist.org/packages/" + pkgname + ".json : ")
	}
	resp, err := http.Get("https://packagist.org/packages/" + pkgname + ".json")
	if err != nil {
		fmt.Printf(" [W] Error when trying to request https://packagist.org/packages/"+pkgname+".json : %s\n", err)
		return false
	}
	defer resp.Body.Close()
	if c.Verbose {
		fmt.Printf("%s\n", resp.Status)
	}
	if resp.StatusCode == http.StatusOK {
		return true
	} else if resp.StatusCode == 429 {
		fmt.Printf(" [!] Server responded with 429 (Too many requests), throttling and retrying...\n")
		time.Sleep(10 * time.Second)
		retry = retry + 1
		return c.isAvailableInPublic(pkgname, retry)
	}
	return false
}

// localReference checks if the package version is in fact a reference to filesystem
func (c *ComposerLookup) localReference(pkgversion string) bool {
	return strings.HasPrefix(strings.ToLower(pkgversion), "file:")
}

// urlReference checks if the package version is in fact a reference to a direct URL
func (c *ComposerLookup) urlReference(pkgversion string) bool {
	pkgversion = strings.ToLower(pkgversion)
	return strings.HasPrefix(pkgversion, "http:") || strings.HasPrefix(pkgversion, "https:")
}

// gitReference checks if the package version is in fact a reference to a remote git repository
func (c *ComposerLookup) gitReference(pkgversion string) bool {
	pkgversion = strings.ToLower(pkgversion)
	gitResources := []string{"git+ssh:", "git+http:", "git+https:", "git:"}
	for _, r := range gitResources {
		if strings.HasPrefix(pkgversion, r) {
			return true
		}
	}
	return false
}

// GetPackageCount returns the number of packages
func (c *ComposerLookup) GetPackageCount() int {
	return len(c.Packages)
}

// GetLanguage returns the language name
func (c *ComposerLookup) GetLanguage() string {
	return "composer"
}

// SetContext sets the context for the resolver
func (c *ComposerLookup) SetContext(ctx context.Context) {
	c.ctx = ctx
}

// SetTimeout sets the timeout for requests
func (c *ComposerLookup) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// SetRateLimit sets the rate limit for requests
func (c *ComposerLookup) SetRateLimit(rate int) {
	c.rateLimit = rate
}

// GetPackageDetails returns detailed information about packages
func (c *ComposerLookup) GetPackageDetails() []types.PackageDetail {
	if len(c.packageDetails) == 0 {
		c.buildPackageDetails()
	}
	return c.packageDetails
}

// buildPackageDetails builds detailed package information
func (c *ComposerLookup) buildPackageDetails() {
	c.packageDetails = []types.PackageDetail{}
	
	for _, pkg := range c.Packages {
		detail := types.PackageDetail{
			Name:    pkg.Name,
			Version: pkg.Version,
			Type:    "dependency",
			Metadata: map[string]interface{}{
				"original_version": pkg.Version,
			},
		}
		
		// Check if package is vulnerable
		if !c.isAvailableInPublic(pkg.Name, 0) {
			detail.Vulnerable = true
			detail.Reason = "Package not available in public Packagist registry"
		}
		
		c.packageDetails = append(c.packageDetails, detail)
	}
}
