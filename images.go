package latest

import (
	"context"
	"fmt"
	"github.com/csmith/latest/internal"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"strings"
)

type ImageOption = func(*imageOptions)

// WithContainerRegistry specifies the default registry to use if image names aren't
// fully-qualified.
//
// If not set, defaults to "docker.io".
func WithContainerRegistry(registry string) ImageOption {
	return func(options *imageOptions) {
		options.registry = registry
	}
}

// WithContainerAuth specifies auth credentials to use when talking to the image
// registry.
//
// If not set, will attempt to use the docker config file.
func WithContainerAuth(username string, password string) ImageOption {
	return func(options *imageOptions) {
		options.username = username
		options.password = password
	}
}

type imageOptions struct {
	registry string
	username string
	password string
}

func (o *imageOptions) craneAuth() crane.Option {
	var authOpt crane.Option

	if o.username == "" || o.password == "" {
		authOpt = crane.WithAuthFromKeychain(authn.DefaultKeychain)
	} else {
		authOpt = crane.WithAuth(&authn.Basic{
			Username: o.username,
			Password: o.password,
		})
	}

	return authOpt
}

func (o *imageOptions) imageName(ref string) string {
	if index := strings.IndexByte(ref, '.'); index != -1 && index < strings.IndexByte(ref, '/') {
		return ref
	} else if o.registry != "" {
		return fmt.Sprintf("%s/%s", o.registry, ref)
	} else {
		return fmt.Sprintf("docker.io/%s", ref)
	}
}

// ImageDigest queries an image registry and returns the latest digest of an
// image with the given name.
func ImageDigest(ctx context.Context, name string, options ...ImageOption) (string, error) {
	o := internal.ResolveOptions(options)
	return crane.Digest(o.imageName(name), o.craneAuth(), crane.WithContext(ctx))
}

type ImageTagOption = func(*imageTagOptions)

type imageTagOptions struct {
	imageOptions
	tagOptions
}

func WithImageOptions(imageOptions ...ImageOption) ImageTagOption {
	return func(options *imageTagOptions) {
		for i := range imageOptions {
			imageOptions[i](&options.imageOptions)
		}
	}
}

func WithTagOptions(tagOptions ...TagOption) ImageTagOption {
	return func(options *imageTagOptions) {
		for i := range tagOptions {
			tagOptions[i](&options.tagOptions)
		}
	}
}

// ImageTag queries an image registry for the available tags for the given
// image, and returns the latest semver tag.
func ImageTag(ctx context.Context, name string, options ...ImageTagOption) (string, error) {
	o := internal.ResolveOptions(options)
	tags, err := crane.ListTags(o.imageName(name), o.craneAuth(), crane.WithContext(ctx))
	if err != nil {
		return "", err
	}

	return o.tagOptions.latest(tags)
}
