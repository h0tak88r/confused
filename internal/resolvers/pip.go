package resolvers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/h0tak88r/confused/internal/types"
)

// PythonLookup represents a collection of python packages to be tested for dependency confusion.
type PythonLookup struct {
	Packages       []string
	Verbose        bool
	ctx            context.Context
	timeout        time.Duration
	rateLimit      int
	packageDetails []types.PackageDetail
}

// NewPythonLookup constructs a `PythonLookup` struct and returns it
func NewPythonLookup(verbose bool) types.PackageResolver {
	return &PythonLookup{
		Packages:       []string{},
		Verbose:        verbose,
		ctx:            context.Background(),
		timeout:        30 * time.Second,
		rateLimit:      100,
		packageDetails: []types.PackageDetail{},
	}
}

// ReadPackagesFromFile reads package information from a python `requirements.txt` file
//
// Returns any errors encountered
func (p *PythonLookup) ReadPackagesFromFile(filename string) error {
	rawfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	line := ""
	for _, l := range strings.Split(string(rawfile), "\n") {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "#") {
			continue
		}
		if len(l) > 0 {
			// Support line continuation
			if strings.HasSuffix(l, "\\") {
				line += l[:len(l) - 1]
				continue
			}
			line += l
			pkgrow := strings.FieldsFunc(line, p.pipSplit)
			if len(pkgrow) > 0 {
				p.Packages = append(p.Packages, strings.TrimSpace(pkgrow[0]))
			}
			// reset the line variable
			line = ""
		}
	}
	return nil
}

// PackagesNotInPublic determines if a python package does not exist in the pypi package repository.
//
// Returns a slice of strings with any python packages not in the pypi package repository
func (p *PythonLookup) PackagesNotInPublic() []string {
	notavail := []string{}
	for _, pkg := range p.Packages {
		if !p.isAvailableInPublic(pkg) {
			notavail = append(notavail, pkg)
		}
	}
	return notavail
}

func (p *PythonLookup) pipSplit(r rune) bool {
	delims := []rune{
		'=',
		'<',
		'>',
		'!',
		' ',
		'~',
		'#',
		'[',
	}
	return inSlice(r, delims)
}

// isAvailableInPublic determines if a python package exists in the pypi package repository.
//
// Returns true if the package exists in the pypi package repository.
func (p *PythonLookup) isAvailableInPublic(pkgname string) bool {
	if p.Verbose {
		fmt.Print("Checking: https://pypi.org/project/" + pkgname + "/ : ")
	}
	resp, err := http.Get("https://pypi.org/project/" + pkgname + "/")
	if err != nil {
		fmt.Printf(" [W] Error when trying to request https://pypi.org/project/"+pkgname+"/ : %s\n", err)
		return false
	}
	if p.Verbose {
		fmt.Printf("%s\n", resp.Status)
	}
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}

// GetPackageCount returns the number of packages
func (p *PythonLookup) GetPackageCount() int {
	return len(p.Packages)
}

// GetLanguage returns the language name
func (p *PythonLookup) GetLanguage() string {
	return "pip"
}

// SetContext sets the context for the resolver
func (p *PythonLookup) SetContext(ctx context.Context) {
	p.ctx = ctx
}

// SetTimeout sets the timeout for requests
func (p *PythonLookup) SetTimeout(timeout time.Duration) {
	p.timeout = timeout
}

// SetRateLimit sets the rate limit for requests
func (p *PythonLookup) SetRateLimit(rate int) {
	p.rateLimit = rate
}

// GetPackageDetails returns detailed information about packages
func (p *PythonLookup) GetPackageDetails() []types.PackageDetail {
	if len(p.packageDetails) == 0 {
		p.buildPackageDetails()
	}
	return p.packageDetails
}

// buildPackageDetails builds detailed package information
func (p *PythonLookup) buildPackageDetails() {
	p.packageDetails = []types.PackageDetail{}
	
	for _, pkgName := range p.Packages {
		detail := types.PackageDetail{
			Name:    pkgName,
			Version: "",
			Type:    "dependency",
			Metadata: map[string]interface{}{
				"original_name": pkgName,
			},
		}
		
		// Check if package is vulnerable
		if !p.isAvailableInPublic(pkgName) {
			detail.Vulnerable = true
			detail.Reason = "Package not available in public PyPI registry"
		}
		
		p.packageDetails = append(p.packageDetails, detail)
	}
}