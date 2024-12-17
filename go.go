package latest

import (
	"context"
	"fmt"
	"github.com/csmith/latest/internal"
	"github.com/hashicorp/go-version"
	"net/url"
	"strings"
)

type GoOption = func(*goOption)

type goOption struct {
	os   string
	arch string
	kind string
}

// WithOS indicates that GoRelease should return only files for the specified OS.
func WithOS(os string) GoOption {
	return func(o *goOption) {
		o.os = os
	}
}

// WithArch indicates that GoRelease should return only files for the specified architecture.
func WithArch(arch string) GoOption {
	return func(o *goOption) {
		o.arch = arch
	}
}

// WithKind indicates that GoRelease should return only files of the specified kind.
//
// Valid kinds at time of writing are: "source", "archive" and "installer".
//
// If not specified defaults to "source".
func WithKind(kind string) GoOption {
	return func(o *goOption) {
		o.kind = kind
	}
}

// GoRelease finds the latest release of Go, returning the version, download URL
// and file checksum.
func GoRelease(ctx context.Context, options ...GoOption) (latestVersion, downloadUrl, downloadChecksum string, err error) {
	o := internal.ResolveOptions(options)

	const (
		baseUrl      = "https://golang.org/dl/"
		jsonReleases = baseUrl + "?mode=json"
		defaultKind  = "source"
	)

	if o.kind == "" {
		o.kind = defaultKind
	}

	var releases []struct {
		Version string `json:"version"`
		Files   []struct {
			Filename string `json:"filename"`
			Checksum string `json:"sha256"`
			Kind     string `json:"kind"`
			Os       string `json:"os"`
			Arch     string `json:"arch"`
		} `yaml:"files"`
	}

	if err := internal.FetchJson(ctx, jsonReleases, &releases); err != nil {
		return "", "", "", err
	}

	best := version.Must(version.NewVersion("0.0.0"))
	for i := range releases {
		r := releases[i]
		v, err := version.NewVersion(strings.TrimPrefix(r.Version, "go"))
		if err != nil {
			return "", "", "", fmt.Errorf("unable to parse version '%s': %w", r.Version, err)
		}

		if v.GreaterThan(best) || (v.Equal(best) && strings.Compare(v.Original(), best.Original()) < 0) {
			best = v

			for j := range r.Files {
				if r.Files[j].Kind == o.kind && (o.arch == "" || r.Files[j].Arch == o.arch) && (o.os == "" || r.Files[j].Os == o.os) {
					latestVersion = r.Version
					downloadUrl, _ = url.JoinPath(baseUrl, r.Files[j].Filename)
					downloadChecksum = r.Files[j].Checksum
				}
			}
		}
	}

	if latestVersion == "" {
		return "", "", "", fmt.Errorf("unable to find matching golang release")
	}

	return
}
