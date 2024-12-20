package latest

import (
	"context"
	"fmt"
	"github.com/csmith/apkutils/v2"
	"github.com/csmith/apkutils/v2/keys"
	"github.com/csmith/latest/internal"
	"net/http"
	"net/url"
	"strings"
)

const defaultAlpineMirror = "https://dl-cdn.alpinelinux.org/alpine/"
const defaultAlpineArch = "x86_64"

// AlpineReleaseOptions defines options for calls to AlpineRelease
type AlpineReleaseOptions struct {
	// The alpine mirror to use. Defaults to `https://dl-cdn.alpinelinux.org/alpine/`
	Mirror string
	// The architecture to select. Defaults to "x86_64".
	Arch string
	// The flavour of alpine to return download links and checksums for. Defaults to "standard".
	Flavour string
}

// AlpineRelease finds the latest release of Alpine Linux, returning the version, download URL
// and file checksum.
func AlpineRelease(ctx context.Context, options *AlpineReleaseOptions) (latestVersion, downloadUrl, downloadChecksum string, err error) {
	o := internal.ApplyDefaults(
		&AlpineReleaseOptions{
			Mirror:  defaultAlpineMirror,
			Arch:    defaultAlpineArch,
			Flavour: "standard",
		},
		options,
	)

	base, err := url.JoinPath(o.Mirror, "latest-stable", "releases", o.Arch)
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
		if strings.TrimPrefix(release.Flavour, "alpine-") == o.Flavour {
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

// AlpinePackageCache facilitates storing information about available Alpine
// packages between calls of AlpinePackage.
type AlpinePackageCache interface {
	// Put stores the given data in the cache
	Put(map[string]*apkutils.PackageInfo) error
	// Get retrieves the existing cached data, or an empty map if not available
	Get() (map[string]*apkutils.PackageInfo, error)
}

type inMemoryAlpinePackage struct {
	packages map[string]*apkutils.PackageInfo
}

func (i *inMemoryAlpinePackage) Put(m map[string]*apkutils.PackageInfo) error {
	i.packages = m
	return nil
}

func (i *inMemoryAlpinePackage) Get() (map[string]*apkutils.PackageInfo, error) {
	return i.packages, nil
}

// NewInMemoryAlpinePackageCache creates a new AlpinePackageCache that simply
// stores the package data in memory.
func NewInMemoryAlpinePackageCache() AlpinePackageCache {
	return &inMemoryAlpinePackage{
		packages: make(map[string]*apkutils.PackageInfo),
	}
}

type AlpinePackageOptions struct {
	// The cache to use for persisting package information. Defaults to a new
	// in-memory cache, which will cause the package list to be updated every
	// time AlpinePackage is called.
	Cache AlpinePackageCache
	// The provider for Alpine Linux public keys, used to verify the signature of the APKINDEX.
	// Defaults to using all known keys provided by apkutil.
	KeyProvider apkutils.KeyProvider
	// The alpine mirror to use. Defaults to `https://dl-cdn.alpinelinux.org/alpine/`
	Mirror string
	// The architecture to use. Defaults to "x86_64".
	Arch string
	// The alpine branch to get packages for. Defaults to "latest-stable"
	Branch string
	// The alpine repository to get packages from. Defaults to "main"
	Repository string
}

// AlpinePackage retrieves the latest version of the given package, the url it can be downloaded from, and the names
// of its declared dependencies. It does not return a checksum as the APK files themselves can be verified using
// their signatures.
//
// Knowing the latest version of a package involves downloading and parsing the APKINDEX file, which is a moderately
// expensive operation. If you are calling this func multiple times, you should populate the Cache option with a
// cache that will maintain state. See [NewInMemoryAlpinePackageCache] for a simple in-memory cache that will suffice
// for most purposes.
//
// If calling this func with multiple different options (such as repository, branch or architecture), care should be
// given to ensure that a unique cache is provided for each combination.
//
// pkg can be a package name, or any token that is provided by a package, e.g. "so:libssl.so.3" or "cmd:busybox".
func AlpinePackage(ctx context.Context, pkg string, options *AlpinePackageOptions) (latestVersion, downloadUrl string, dependencies []string, err error) {
	o := internal.ApplyDefaults(&AlpinePackageOptions{
		Cache:       NewInMemoryAlpinePackageCache(),
		KeyProvider: keys.All,
		Mirror:      defaultAlpineMirror,
		Arch:        defaultAlpineArch,
		Branch:      "latest-stable",
		Repository:  "main",
	}, options)

	packages, err := o.Cache.Get()
	if err != nil {
		return "", "", nil, err
	}

	base, err := url.JoinPath(o.Mirror, o.Branch, o.Repository, o.Arch)
	if err != nil {
		return "", "", nil, err
	}

	if len(packages) == 0 {
		packages, err = (func() (map[string]*apkutils.PackageInfo, error) {
			apkIndex, err := url.JoinPath(base, "APKINDEX.tar.gz")
			if err != nil {
				return nil, err
			}

			println(apkIndex)

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, apkIndex, nil)
			if err != nil {
				return nil, err
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return nil, err
			}

			defer res.Body.Close()
			return apkutils.ReadApkIndex(res.Body, o.KeyProvider)
		})()
		if err != nil {
			return "", "", nil, err
		}

		err = o.Cache.Put(packages)
		if err != nil {
			return "", "", nil, err
		}
	}

	p := packages[pkg]
	if p == nil {
		return "", "", nil, fmt.Errorf("couldn't find package %s", pkg)
	}

	latestVersion = p.Version
	downloadUrl, err = url.JoinPath(base, fmt.Sprintf("%s-%s.apk", p.Name, p.Version))
	if err != nil {
		return "", "", nil, err
	}
	dependencies = p.Dependencies
	return
}
