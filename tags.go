package latest

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"slices"
	"strings"

	"github.com/hashicorp/go-version"
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
	// Optional function to apply to tags once prefixes/suffixes have been
	// trimmed. Can be used to normalize tags that don't follow semver
	// conventions, e.g. replacing underscores with dots in curl's
	// `curl-8_13_0` tags.
	Transform func(string) string
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
		v, err := t.parse(tags[i])
		if err != nil {
			return "", err
		}

		if v == nil || !t.allows(v) {
			continue
		}

		if v.GreaterThan(best) || (v.Equal(best) && strings.Compare(v.Original(), best.Original()) < 0) {
			best = v
			bestRaw = tags[i]
		}
	}

	if bestRaw == "" {
		return "", errors.New("no valid tags found")
	}

	return bestRaw, nil
}

// parse normalises the tag and parses it as a version. It returns a nil version if the tag should be skipped.
func (t *TagOptions) parse(tag string) (*version.Version, error) {
	stripped := tag
	for _, prefix := range t.TrimPrefixes {
		stripped = strings.TrimPrefix(stripped, prefix)
	}
	for _, suffix := range t.TrimSuffixes {
		stripped = strings.TrimSuffix(stripped, suffix)
	}
	if t.Transform != nil {
		stripped = t.Transform(stripped)
	}

	if t.IgnoreDates && dateRegexp.MatchString(stripped) {
		return nil, nil
	}

	v, err := version.NewVersion(stripped)
	if err != nil {
		if t.IgnoreErrors {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to parse tag '%s': %w", stripped, err)
	}

	return v, nil
}

func (t *TagOptions) allows(v *version.Version) bool {
	prerelease := v.Prerelease()
	if prerelease != "" && t.IgnorePreRelease {
		return false
	}

	if len(t.PreReleases) > 0 {
		if prerelease == "" {
			return false
		}

		matches := func(allowed string) bool { return strings.EqualFold(prerelease, allowed) }
		if !slices.ContainsFunc(t.PreReleases, matches) {
			return false
		}
	}

	return t.MajorVersionMax >= v.Segments()[0]
}
