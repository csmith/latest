package main

import (
	"context"
	"fmt"
	"github.com/csmith/apkutils/v2/keys"
	"github.com/csmith/latest"
)

func main() {
	version, url, dependencies, err := latest.AlpinePackage(
		context.Background(),
		"cmake",
		&latest.AlpinePackageOptions{
			AlpinePackagesOptions: latest.AlpinePackagesOptions{
				Arch:        "aarch64",
				Branch:      "edge",
				KeyProvider: keys.Aarch64,
			},
		},
	)
	if err != nil {
		panic(err)
	}

	// e.g. "3.31.3-r0 https://dl-cdn.alpinelinux.org/alpine/edge/main/aarch64/cmake-3.31.3-r0.apk [so:libarchive.so.13 so:libc.musl-aarch64.so.1 so:libcrypto.so.3 so:libexpat.so.1 so:libgcc_s.so.1 so:librhash.so.1 so:libssl.so.3 so:libstdc++.so.6 so:libuv.so.1 so:libz.so.1]"
	fmt.Println(version, url, dependencies)
}
