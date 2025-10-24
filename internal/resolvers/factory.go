package resolvers

import (
	"fmt"
	"github.com/h0tak88r/confused2/internal/types"
)

// GetResolverForLanguage returns a resolver for the given language
func GetResolverForLanguage(language string) (types.PackageResolver, error) {
	switch language {
	case "pip":
		return NewPythonLookup(false), nil
	case "npm":
		return NewNPMLookup(false), nil
	case "composer":
		return NewComposerLookup(false), nil
	case "mvn":
		return NewMVNLookup(false), nil
	case "rubygems":
		return NewRubyGemsLookup(false), nil
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
}

// GetResolverForLanguageWithVerbose returns a resolver for the given language with verbose setting
func GetResolverForLanguageWithVerbose(language string, verbose bool) (types.PackageResolver, error) {
	switch language {
	case "pip":
		return NewPythonLookup(verbose), nil
	case "npm":
		return NewNPMLookup(verbose), nil
	case "composer":
		return NewComposerLookup(verbose), nil
	case "mvn":
		return NewMVNLookup(verbose), nil
	case "rubygems":
		return NewRubyGemsLookup(verbose), nil
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
}
