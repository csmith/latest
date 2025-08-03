package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/csmith/apkutils/v2"
	"github.com/csmith/apkutils/v2/keys"
	"github.com/csmith/envflag/v2"
	"github.com/csmith/latest"
)

var showJson = flag.Bool("json", false, "Provide output in json")
var arch = flag.String("arch", "x86_64", "Architecture to query")
var branch = flag.String("branch", "latest-stable", "Alpine branch to query")

func main() {
	envflag.Parse()

	var keyProvider apkutils.KeyProvider
	switch *arch {
	case "aarch64":
		keyProvider = keys.Aarch64
	case "x86_64":
		keyProvider = keys.X86_64
	default:
		fmt.Printf("Unknown/unsupported architecture specified.")
		os.Exit(1)
	}

	packages, err := latest.AlpinePackages(
		context.Background(),
		&latest.AlpinePackagesOptions{
			Arch:        *arch,
			Branch:      *branch,
			KeyProvider: keyProvider,
		},
	)
	if err != nil {
		fmt.Printf("Error getting alpine packages: %s\n", err.Error())
		os.Exit(1)
	}

	if *showJson {
		bytes, err := json.Marshal(packages)
		if err != nil {
			fmt.Printf("Error marshaling json: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(bytes))
		return
	}

	fmt.Printf("%d packages found\n", len(packages))
	for name, pkg := range packages {
		fmt.Printf("%s: %s\n", name, pkg.Version)
	}
}
