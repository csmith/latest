package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/csmith/envflag/v2"
	"github.com/csmith/latest"
)

var showJson = flag.Bool("json", false, "Provide output in json")

func main() {
	envflag.Parse()
	var imageName string
	if flag.NArg() < 1 {
		fmt.Printf("Must specify which container image\n")
		os.Exit(1)
	}
	imageName = flag.Arg(0)

	digest, err := latest.ImageDigest(
		context.Background(),
		imageName,
		nil,
	)
	if err != nil {
		fmt.Printf("Error getting image digest: %s\n", err.Error())
		os.Exit(1)
	}

	details := struct {
		Image  string `json:"image"`
		Digest string `json:"digest"`
	}{
		Image:  imageName,
		Digest: digest,
	}

	if *showJson {
		bytes, err := json.Marshal(&details)
		if err != nil {
			fmt.Printf("Error marshaling json: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(bytes))
		return
	}

	fmt.Printf("Image: %s\nDigest: %s\n", imageName, digest)
}
