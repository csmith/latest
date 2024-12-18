package main

import (
	"context"
	"fmt"
	"github.com/csmith/latest"
)

func main() {
	version, url, checksum, err := latest.PostgresRelease(context.Background(), latest.WithMaximumMajorVersion(16))
	if err != nil {
		panic(err)
	}

	// e.g. "16.6 https://ftp.postgresql.org/pub/source/v16.6/postgresql-16.6.tar.bz2 23369cdaccd45270ac5dcc30fa9da205d5be33fa505e1f17a0418d2caeca477b"
	fmt.Println(version, url, checksum)
}
