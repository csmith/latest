package latest

import (
	"context"
	"maps"
	"slices"
	"strings"

	"github.com/csmith/gitrefs"
	"github.com/csmith/latest/internal"
)

const gitTagPrefix = "refs/tags/"

type GitTagOptions struct {
	TagOptions
	Username string
	Password string
}

var defaultGitTagOptions = GitTagOptions{
	TagOptions: defaultTagOptions,
}

// GitTag lists the tags available in the specified repository, and returns the
// latest semver tag.
//
// The repository must be specified as an HTTP/HTTPS url.
func GitTag(ctx context.Context, repo string, options *GitTagOptions) (string, error) {
	o := internal.ApplyDefaults(&defaultGitTagOptions, options)
	o.TrimPrefixes = append([]string{gitTagPrefix}, o.TrimPrefixes...)

	gitrefsOptions := []gitrefs.Option{
		gitrefs.TagsOnly(),
		gitrefs.WithContext(ctx),
	}
	if options.Username != "" || options.Password != "" {
		gitrefsOptions = append(gitrefsOptions, gitrefs.WithAuth(options.Username, options.Password))
	}
	refs, err := gitrefs.Fetch(repo, gitrefsOptions...)
	if err != nil {
		return "", err
	}

	tag, err := o.latest(slices.Collect(maps.Keys(refs)))
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(tag, gitTagPrefix), nil
}
