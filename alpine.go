package latest

import (
	"context"
	"fmt"
	"github.com/csmith/latest/internal"
	"net/url"
	"strings"
)

type AlpineReleaseOption = func(options *alpineReleaseOptions)

type alpineReleaseOptions struct {
	mirror  string
	arch    string
	flavour string
}

// WithAlpineMirror indicates which alpine mirror should be used to retrieve
// alpine package information. Defaults to https://dl-cdn.alpinelinux.org/alpine/.
func WithAlpineMirror(mirror string) AlpineReleaseOption {
	return func(o *alpineReleaseOptions) {
		o.mirror = mirror
	}
}

// WithAlpineArch specifies the architecture to be queried. Defaults to x86_64.
func WithAlpineArch(arch string) AlpineReleaseOption {
	return func(o *alpineReleaseOptions) {
		o.arch = arch
	}
}

// WithAlpineFlavour specifies the flavour of Alpine to be queried. Defaults to standard.
func WithAlpineFlavour(flavour string) AlpineReleaseOption {
	return func(o *alpineReleaseOptions) {
		o.flavour = flavour
	}
}

// AlpineRelease finds the latest release of Alpine Linux, returning the version, download URL
// and file checksum.
func AlpineRelease(ctx context.Context, options ...AlpineReleaseOption) (latestVersion, downloadUrl, downloadChecksum string, err error) {
	o := internal.ResolveOptionsWithDefaults(options, &alpineReleaseOptions{
		mirror:  "https://dl-cdn.alpinelinux.org/alpine/",
		arch:    "x86_64",
		flavour: "standard",
	})

	base, err := url.JoinPath(o.mirror, "latest-stable", "releases", o.arch)
	if err != nil {
		return "", "", "", err
	}

	yamlUrl, err := url.JoinPath(base, "latest-releases.yaml")
	if err != nil {
		return "", "", "", err
	}

	var releases []struct {
		Flavour  string `yaml:"flavor"`
		File     string `yaml:"file"`
		Checksum string `yaml:"sha256"`
		Version  string `yaml:"version"`
	}

	if err := internal.FetchYaml(ctx, yamlUrl, &releases); err != nil {
		return "", "", "", err
	}

	for _, release := range releases {
		if strings.TrimPrefix(release.Flavour, "alpine-") == o.flavour {
			latestVersion = release.Version
			downloadChecksum = release.Checksum
			downloadUrl, err = url.JoinPath(base, release.File)

			if err != nil {
				return "", "", "", err
			}

			return
		}
	}

	return "", "", "", fmt.Errorf("couldn't find matching alpine release")
}
