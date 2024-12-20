package main

import (
	"context"
	"fmt"
	"github.com/csmith/latest"
)

func main() {
	version, url, checksum, err := latest.AlpineRelease(
		context.Background(),
		&latest.AlpineReleaseOptions{
			Arch:    "aarch64",
			Flavour: "minirootfs",
		},
	)
	if err != nil {
		panic(err)
	}

	// e.g. "3.21.0 https://dl-cdn.alpinelinux.org/alpine/latest-stable/releases/aarch64/alpine-minirootfs-3.21.0-aarch64.tar.gz f31202c4070c4ef7de9e157e1bd01cb4da3a2150035d74ea5372c5e86f1efac1"
	fmt.Println(version, url, checksum)
}
