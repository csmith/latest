package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/csmith/envflag/v2"
	"github.com/csmith/latest/v3"
)

var showJson = flag.Bool("json", false, "Provide output in json")

func main() {
	envflag.Parse()
	var containerName string
	if flag.NArg() < 1 {
		_, _ = fmt.Fprintf(os.Stderr, "Must specify which container name")
		os.Exit(1)
	} else {
		containerName = flag.Arg(0)
	}
	version, err := latest.ImageTag(context.Background(), containerName, nil)
	if err != nil {
		fmt.Println("Error getting latest container release: " + err.Error())
	}
	details := struct {
		Version string `json:"version"`
	}{
		Version: version,
	}
	if *showJson {
		err = nil
		bytes, err := json.Marshal(&details)
		if err != nil {
			fmt.Println("Error printing json: " + err.Error())
		}
		fmt.Println(string(bytes))
		return
	}
	fmt.Printf("Version: %s\n", details.Version)
}
