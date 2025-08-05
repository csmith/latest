package latest

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"math"
	"regexp"
	"strings"
)

// TagOptions defines options for functions that need to deal with semver tags.
type TagOptions struct {
	// Should tags that look like dates (starting yyyymmdd or yyyy-mm-dd) be ignored?
	IgnoreDates bool
	// Should tags that are not parseable as semver be silently ignored?
	IgnoreErrors bool
	// Should pre-releases (-alpha, -beta, etc) be ignored?
	IgnorePreRelease bool
	// Only consider tags with prerelease identifiers matching one of these (case-insensitive)
	PreReleases []string
	// Strings to remove from the start of tags. Processed in order.
	TrimPrefixes []string
	// Strings to remove from the ends of tags. Processed in order.
	TrimSuffixes []string
	// The maximum major version to consider.
	MajorVersionMax int
}

var defaultTagOptions = TagOptions{
	MajorVersionMax: math.MaxInt,
}

var dateRegexp = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})|(\d{8})`)

func (t *TagOptions) latest(tags []string) (string, error) {
	best := version.Must(version.NewVersion("0.0.0"))
	bestRaw := ""

	for i := range tags {
		stripped := tags[i]
		for _, prefix := range t.TrimPrefixes {
			stripped = strings.TrimPrefix(stripped, prefix)
		}
		for _, suffix := range t.TrimSuffixes {
			stripped = strings.TrimSuffix(stripped, suffix)
		}

		if t.IgnoreDates && dateRegexp.MatchString(stripped) {
			continue
		}

		v, err := version.NewVersion(stripped)
		if err != nil {
			if t.IgnoreErrors {
				continue
			} else {
				return "", fmt.Errorf("unable to parse tag '%s': %w", stripped, err)
			}
		}

		if v.Prerelease() != "" && t.IgnorePreRelease {
			continue
		}

		if len(t.PreReleases) > 0 {
			prerelease := v.Prerelease()
			if prerelease == "" {
				continue
			}

			found := false
			for _, allowed := range t.PreReleases {
				if strings.EqualFold(prerelease, allowed) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if t.MajorVersionMax < v.Segments()[0] {
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
