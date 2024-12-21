package main

import (
	"context"
	"fmt"
	"github.com/csmith/apkutils/v2/keys"
	"github.com/csmith/latest"
)

func main() {
	packages, err := latest.AlpinePackages(
		context.Background(),
		&latest.AlpinePackagesOptions{
			Arch:        "aarch64",
			Branch:      "edge",
			KeyProvider: keys.Aarch64,
		},
	)
	if err != nil {
		panic(err)
	}

	// e.g. "13381 packages/providers found. Cmake version = 3.31.3-r0"
	fmt.Printf("%d packages/providers found. Cmake version = %s\n", len(packages), packages["cmake"].Version)
}
