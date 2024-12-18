package latest

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"regexp"
	"strings"
)

type TagOption = func(*tagOptions)

type tagOptions struct {
	ignoreDates      bool
	ignoreErrors     bool
	ignorePreRelease bool
	trimPrefixes     []string
	trimSuffixes     []string
	majorVersionMax  *int
}

// WithIgnoreDates indicates that tags that look like dates (starting yyyymmdd or yyyy-mm-dd)
// should be ignored.
func WithIgnoreDates() TagOption {
	return func(o *tagOptions) {
		o.ignoreDates = true
	}
}

// WithIgnoreErrors indicates tags that are not parseable as semver should be silently ignored.
func WithIgnoreErrors() TagOption {
	return func(o *tagOptions) {
		o.ignoreErrors = true
	}
}

// WithTrimPrefix indicates that if the given prefix is found on a tag, it should be ignored.
// This option may be specified multiple times. Prefixes will be stripped in order.
func WithTrimPrefix(prefix string) TagOption {
	return func(o *tagOptions) {
		o.trimPrefixes = append(o.trimPrefixes, prefix)
	}
}

// WithTrimSuffix indicates that if the given suffix is found on a tag, it should be ignored.
// This option may be specified multiple times. Suffixes will be stripped in order.
func WithTrimSuffix(suffix string) TagOption {
	return func(o *tagOptions) {
		o.trimSuffixes = append(o.trimSuffixes, suffix)
	}
}

// WithMaximumMajorVersion constraints the result to have a major version not
// more than the one specified.
func WithMaximumMajorVersion(version int) TagOption {
	return func(o *tagOptions) {
		o.majorVersionMax = &version
	}
}

var dateRegexp = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})|(\d{8})`)

func (t *tagOptions) latest(tags []string) (string, error) {
	best := version.Must(version.NewVersion("0.0.0"))
	bestRaw := ""

	for i := range tags {
		stripped := tags[i]
		for _, prefix := range t.trimPrefixes {
			stripped = strings.TrimPrefix(stripped, prefix)
		}
		for _, suffix := range t.trimSuffixes {
			stripped = strings.TrimSuffix(stripped, suffix)
		}

		if t.ignoreDates && dateRegexp.MatchString(stripped) {
			continue
		}

		v, err := version.NewVersion(stripped)
		if err != nil {
			if t.ignoreErrors {
				continue
			} else {
				return "", fmt.Errorf("unable to parse tag '%s': %w", stripped, err)
			}
		}

		if v.Prerelease() != "" && t.ignorePreRelease {
			continue
		}

		if t.majorVersionMax != nil && *t.majorVersionMax < v.Segments()[0] {
			continue
		}

		if v.GreaterThan(best) || (v.Equal(best) && strings.Compare(v.Original(), best.Original()) < 0) {
			best = v
			bestRaw = tags[i]
		}
	}

	if bestRaw == "" {
		return "", fmt.Errorf("no valid tags found")
	}

	return bestRaw, nil
}
