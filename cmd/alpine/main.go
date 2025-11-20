package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/csmith/envflag/v2"
	"github.com/csmith/latest/v2"
)

var showJson = flag.Bool("json", false, "Provide output in json")

func main() {
	envflag.Parse()
	version, download, checksum, err := latest.AlpineRelease(context.Background(), nil)
	if err != nil {
		fmt.Println("Error getting latest alpine release: " + err.Error())
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
