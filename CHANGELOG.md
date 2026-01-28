# Changelog

## Unreleased

### Other changes

- Migrated to github.com/goccy/go-yaml for YAML parsing

## 3.0.0 - 2025-12-22

### Breaking changes

- `GitTag` now also returns the commit ID the tag points to.

### Bug fixes

- Fix panic in `GitTag` if no options are provided.

## 2.0.1 - 2025-11-20

- Corrected module path to github.com/csmith/latest/v2

## 2.0.0 - 2025-11-20

### Breaking changes

- `GitTag` now takes a `GitTagOptions` struct that wraps `TagOptions`.

### Features

- `GitTag` can now pass on a username and password to authenticate with the
  git server.

## 1.4.0 - 2025-11-06

### Bug fixes

- Fix getting latest git version in repositories with dereferenced
  tags (thanks @Greboid)

## 1.3.1 - 2025-08-25

### Other changes

- Minor dependency updates
- Minimum go version is now 1.24 instead of an arbitrary patch release

## 1.3.0 - 2025-08-05

### Other changes

- TagOptions now includes a `PreReleases` field. If set, only tags with a
  pre-release contained in the slice will be selected. This is useful for
  selecting beta/RCs specifically, or convoluted versioning schemes with
  components in like "v1.2.3-frontend".

### Bug fixes

- The `container` executable now checks for the latest image tag, instead of
  incorrectly checking for alpine packages.

## 1.2.0 - 2025-08-03

### Other changes

- Added executable utilities to exercise the library functions.
  Thanks @Greboid.

## 1.1.1 - 2025-07-16

### Bug fixes 

- Fixed panic if passing `nil` to use default alpine package options.
  Thanks @Greboid.

## 1.1.0 - 2025-07-09

### Other changes

- Minor dependency updates

## 1.0.0 - 2024-12-21

_Initial release._
