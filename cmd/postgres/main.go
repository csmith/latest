package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/csmith/envflag/v2"
	"github.com/csmith/latest/v2"
)

var showJson = flag.Bool("json", false, "Provide output in json")

func main() {
	envflag.Parse()
	options := latest.TagOptions{}
	options.IgnorePreRelease = true
	if flag.NArg() >= 1 {
		majorVersion, err := strconv.Atoi(flag.Arg(0))
		if err != nil {
			os.Exit(1)
		}
		options.MajorVersionMax = majorVersion
	}
	version, download, checksum, err := latest.PostgresRelease(context.Background(), &options)
	if err != nil {
		fmt.Println("Error getting the latest postgres version: " + err.Error())
	}
	details := struct {
		Version  string `json:"version"`
		Download string `json:"download"`
		Checksum string `json:"checksum"`
	}{
		Version:  version,
		Download: download,
		Checksum: checksum,
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
	fmt.Printf("Version: %s\nDownload: %s\nChecksum: %s\n",
		details.Version, details.Download, details.Checksum)
}
