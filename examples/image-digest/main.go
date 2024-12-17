package main

import (
	"context"
	"fmt"
	"github.com/csmith/latest"
)

func main() {
	digest, err := latest.ImageDigest(context.Background(), "alpine:latest")
	if err != nil {
		panic(err)
	}

	// e.g. "sha256:21dc6063fd678b478f57c0e13f47560d0ea4eeba26dfc947b2a4f81f686b9f45"
	fmt.Println(digest)
}
