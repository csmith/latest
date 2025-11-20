package latest

import (
	"context"
	"fmt"
	"strings"

	"github.com/csmith/latest/v2/internal"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
)

// ImageOptions defines options for image-related functions.
type ImageOptions struct {
	// The registry to use if image names aren't fully-qualified.
	// Defaults to "docker.io".
	Registry string
	// The username to use to authenticate to registries.
	// If not set, will attempt to use the docker config file.
	Username string
	// The password to use to authenticate to registries.
	// If not set, will attempt to use the docker config file.
	Password string
}

func (o *ImageOptions) craneAuth() crane.Option {
	var authOpt crane.Option

	if o.Username == "" || o.Password == "" {
		authOpt = crane.WithAuthFromKeychain(authn.DefaultKeychain)
	} else {
		authOpt = crane.WithAuth(&authn.Basic{
			Username: o.Username,
			Password: o.Password,
		})
	}

	return authOpt
}

func (o *ImageOptions) imageName(ref string) string {
	if index := strings.IndexByte(ref, '.'); index != -1 && index < strings.IndexByte(ref, '/') {
		return ref
	} else {
		return fmt.Sprintf("%s/%s", o.Registry, ref)
	}
}

// ImageDigest queries an image registry and returns the latest digest of an
// image with the given name.
func ImageDigest(ctx context.Context, name string, options *ImageOptions) (string, error) {
	o := internal.ApplyDefaults(
		&ImageOptions{
			Registry: "docker.io",
		},
		options,
	)
	return crane.Digest(o.imageName(name), o.craneAuth(), crane.WithContext(ctx))
}

// ImageTagOptions defines the options for calls to ImageTag.
type ImageTagOptions struct {
	ImageOptions
	TagOptions
}

// ImageTag queries an image registry for the available tags for the given
// image, and returns the latest semver tag.
func ImageTag(ctx context.Context, name string, options *ImageTagOptions) (string, error) {
	o := internal.ApplyDefaults(
		&ImageTagOptions{
			ImageOptions: ImageOptions{
				Registry: "docker.io",
			},
			TagOptions: defaultTagOptions,
		},
		options,
	)

	tags, err := crane.ListTags(o.imageName(name), o.craneAuth(), crane.WithContext(ctx))
	if err != nil {
		return "", err
	}

	return o.TagOptions.latest(tags)
}
