package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mholt/archiver/v3"
	gh "github.com/richardltc/bw_manager/github"
	rjmInternet "github.com/richardltc/bw_manager/rjminternet"
)

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir()
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	// This covers other errors (permissions, etc.)
	return false
}

func main() {
	const green = "\x1b[32m"
	const boldGreen = "\x1b[1;32m"
	const reset = "\x1b[0m"

	// Mocking an app struct for context
	app := struct {
		name    string
		version string
	}{
		name:    "BoxWallet Manager",
		version: "0.0.1",
	}

	// Immediate log
	fmt.Printf("%s%s%s v%s%s%s starting...\n",
		green, app.name, reset,
		green, app.version, reset,
	)

	// Defer runs when the surrounding function exits
	defer fmt.Printf("%s%s%s finished!\n",
		green, app.name, reset,
	)

	fmt.Printf("Checking for updates...\n")

	// Call the function using the package name prefix
	latest_version, downloadUrl, filename, err := gh.GetLatestDownloadinfo()
	if err != nil {
		log.Fatalf("Failed to fetch release: %v", err)
	}

	if !dirExists(latest_version) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nThe latest version of " + green + "BoxWallet" + reset + " is: " + boldGreen + latest_version + reset + "\n\nWould you like to download it? (y/n): ")
		input, _ := reader.ReadString('\n')
		fmt.Printf("Text entered: %s\n", input)
		if input != "y\n" {
			fmt.Println("Update cancelled.")
			return
		}

		dir := latest_version

		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "unable to create directory: %v", err)
		}

		dest := path.Join(dir, filename)

		// The user has chosen to download the latest version...
		if err := rjmInternet.DownloadFile(dest, downloadUrl); err != nil {
			fmt.Fprintf(os.Stderr, "unable to download file: %v - %v", downloadUrl, err)
		}

		if err := archiver.Unarchive(dest, dir); err != nil {
			fmt.Fprintf(os.Stderr, "unable to unarchive file: %v - %v", dest, err)
		}

	} else {
		// The latest version has already been downloaded, as the directory exists
		fmt.Println("You already have the latest version of " + green + "BoxWallet" + reset)

	}

	// fmt.Printf("Latest release found: %s%s%s\n", green, downloadUrl, reset)
}
