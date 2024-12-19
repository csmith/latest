package latest

import (
	"context"
	"fmt"
	"github.com/csmith/latest/internal"
	"strings"
)

func PostgresRelease(ctx context.Context, options ...TagOption) (latest, url, checksum string, err error) {
	o := internal.ResolveOptionsWithDefaults(options, &tagOptions{
		trimSuffixes: []string{"/"},
	})

	const (
		postgresReleaseIndex = "https://ftp.postgresql.org/pub/source/"
		postgresDownloadUrl  = postgresReleaseIndex + "v%[1]s/postgresql-%[1]s.tar.bz2"
		postgresChecksumUrl  = postgresReleaseIndex + "v%[1]s/postgresql-%[1]s.tar.bz2.sha256"
	)

	versions, err := internal.FetchHtmlAndFind(ctx, postgresReleaseIndex, `a[href*="v"]`)
	if err != nil {
		return "", "", "", fmt.Errorf("couldn't find postgres releases: %w", err)
	}

	latest, err = o.latest(versions)
	if err != nil {
		return "", "", "", err
	}

	latest = strings.TrimPrefix(strings.TrimSuffix(latest, "/"), "v")
	url = fmt.Sprintf(postgresDownloadUrl, latest)
	checksum, err = internal.FetchHash(ctx, fmt.Sprintf(postgresChecksumUrl, latest))
	if err != nil {
		return "", "", "", err
	}

	return
}
