package main

import (
	"context"
	"fmt"
	"github.com/csmith/latest"
)

func main() {
	version, url, checksum, err := latest.GoRelease(context.Background(), latest.WithKind("installer"), latest.WithOS("darwin"))
	if err != nil {
		panic(err)
	}

	// e.g. "go1.23.4 https://golang.org/dl/go1.23.4.darwin-arm64.pkg 19c054eaf40c5fac65b027f7443c01382e493d3c8c42cf8b2504832ebddce037"
	fmt.Println(version, url, checksum)
}
