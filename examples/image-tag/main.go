package main

import (
	"context"
	"fmt"
	"github.com/csmith/latest"
)

func main() {
	digest, err := latest.ImageTag(
		context.Background(),
		"alpine",
		latest.WithTagOptions(
			latest.WithIgnoreErrors(),
			latest.WithIgnoreDates(),
		),
	)
	if err != nil {
		panic(err)
	}

	// e.g. "3.21.0"
	fmt.Println(digest)
}
