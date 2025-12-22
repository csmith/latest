package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/csmith/envflag/v2"
	"github.com/csmith/latest/v3"
)

var showJson = flag.Bool("json", false, "Provide output in json")

func main() {
	envflag.Parse()
	var alpinePackage string
	if flag.NArg() < 1 {
		_, _ = fmt.Fprintf(os.Stderr, "Must specify which alpine package")
		os.Exit(1)
	} else {
		alpinePackage = flag.Arg(0)
	}
	version, download, deps, err := latest.AlpinePackage(context.Background(), alpinePackage, nil)
	if err != nil {
		fmt.Println("Error getting latest alpine release: " + err.Error())
	}
	details := struct {
		Version      string   `json:"version"`
		Download     string   `json:"download"`
		Dependencies []string `json:"deps"`
	}{
		Version:      version,
		Download:     download,
		Dependencies: deps,
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
	fmt.Printf("Version: %s\nDownload: %s\nDependencies: %s\n",
		details.Version, details.Download, strings.Join(details.Dependencies, ", "))
}
