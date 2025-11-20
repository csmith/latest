package latest

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/csmith/latest/v2/internal"
	"github.com/hashicorp/go-version"
)

// GoOptions defines options for calling GoRelease.
type GoOptions struct {
	// If set, only return files for this OS.
	Os string
	// If set, only return files for this Architecture.
	Arch string
	// The kind of file to return. Defaults to "source".
	//
	// Valid kinds at time of writing are: "source", "archive" and "installer".
	Kind string
}

// GoRelease finds the latest release of Go, returning the version, download URL
// and file checksum.
func GoRelease(ctx context.Context, options *GoOptions) (latestVersion, downloadUrl, downloadChecksum string, err error) {
	o := internal.ApplyDefaults(
		&GoOptions{
			Kind: "source",
		},
		options,
	)

	const (
		baseUrl      = "https://golang.org/dl/"
		jsonReleases = baseUrl + "?mode=json"
	)

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
				if r.Files[j].Kind == o.Kind && (o.Arch == "" || r.Files[j].Arch == o.Arch) && (o.Os == "" || r.Files[j].Os == o.Os) {
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
