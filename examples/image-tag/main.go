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
		&latest.ImageTagOptions{
			TagOptions: latest.TagOptions{
				IgnoreErrors: true,
				IgnoreDates:  true,
			},
		},
	)
	if err != nil {
		panic(err)
	}

	// e.g. "3.21.0"
	fmt.Println(digest)
}
