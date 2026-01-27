package main

import (
	"fmt"
	"log"

	gh "github.com/richardltc/bw_manager/github"
	rjmInternet "github.com/richardltc/bw_manager/rjminternet"
)

func main() {
	const green = "\x1b[32m"
	const reset = "\x1b[0m"

	// Mocking an app struct for context
	app := struct {
		name    string
		version string
	}{
		name:    "BoxWallet",
		version: "0.0.1",
	}

	// Immediate log
	fmt.Printf("%s%s%s v%s%s%s starting...\n",
		green, app.name, reset,
		green, app.version, reset,
	)

	// Defer runs when the surrounding function exits
	defer fmt.Printf("%sBoxWallet%s v%s0.01%s finished!\n",
		green, reset,
		green, reset,
	)

	fmt.Printf("Checking for updates...\n")

	// Call the function using the package name prefix
	downloadUrl, err := gh.GetLatestDownloadUri()
	if err != nil {
		log.Fatalf("Failed to fetch release: %v", err)
	}

	if err := rjmInternet.DownloadFile("test.tar.gz", downloadUrl); err != nil {
		fmt.Errorf("unable to download file: %v - %v", downloadUrl, err)
	}

	fmt.Printf("Latest release found: %s%s%s\n", green, downloadUrl, reset)
}
