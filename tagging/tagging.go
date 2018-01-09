package tagging

import (
	"fmt"
	"strings"

	"github.com/coreos/go-semver/semver"
)

// DefaultTags returns a set of default suggested tags.
func DefaultTags(ref string) []string {
	if !strings.HasPrefix(ref, "refs/tags/") {
		return []string{"latest"}
	}

	v := stripTagPrefix(ref)

	version, err := semver.NewVersion(v)

	if err != nil {
		return []string{"latest"}
	}

	if version.PreRelease != "" || version.Metadata != "" {
		return []string{
			version.String(),
		}
	}

	if version.Major == 0 {
		return []string{
			fmt.Sprintf("%d.%d", version.Major, version.Minor),
			fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch),
		}
	}

	return []string{
		fmt.Sprint(version.Major),
		fmt.Sprintf("%d.%d", version.Major, version.Minor),
		fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch),
	}
}

// UseDefaultTag to restrict latest tag for default branch.
func UseDefaultTag(ref, defaultBranch string) bool {
	if strings.HasPrefix(ref, "refs/tags/") {
		return true
	}

	if stripHeadPrefix(ref) == defaultBranch {
		return true
	}

	return false
}

// stripHeadPrefix just strips the ref heads prefix.
func stripHeadPrefix(ref string) string {
	return strings.TrimPrefix(ref, "refs/heads/")
}

// stripTagPrefix just strips the ref tags prefix.
func stripTagPrefix(ref string) string {
	ref = strings.TrimPrefix(ref, "refs/tags/")
	ref = strings.TrimPrefix(ref, "v")

	return ref
}
