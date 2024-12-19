package latest

import (
	"context"
	"github.com/csmith/gitrefs"
	"github.com/csmith/latest/internal"
	"maps"
	"slices"
	"strings"
)

const gitTagPrefix = "refs/tags/"

// GitTag lists the tags available in the specified repository, and returns the
// latest semver tag.
//
// The repository must be specified as an HTTP/HTTPS url.
func GitTag(ctx context.Context, repo string, options ...TagOption) (string, error) {
	o := internal.ResolveOptionsWithDefaults(options, &tagOptions{
		trimPrefixes: []string{gitTagPrefix},
	})

	refs, err := gitrefs.Fetch(repo, gitrefs.TagsOnly(), gitrefs.WithContext(ctx))
	if err != nil {
		return "", err
	}

	tag, err := o.latest(slices.Collect(maps.Keys(refs)))
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(tag, gitTagPrefix), nil
}
