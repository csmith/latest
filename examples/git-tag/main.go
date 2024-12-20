package main

import (
	"context"
	"fmt"
	"github.com/csmith/latest"
)

func main() {
	digest, err := latest.GitTag(
		context.Background(),
		"https://github.com/csmith/gitrefs",
		nil,
	)
	if err != nil {
		panic(err)
	}

	// e.g. "v1.3.0"
	fmt.Println(digest)
}
