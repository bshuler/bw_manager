package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

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

func runBoxWallet(subDir string) {
	// 1. Determine the path to the Elixir executable based on OS
	var executable string
	if runtime.GOOS == "windows" {
		executable = filepath.Join(subDir, "boxwallet", "bin", "boxwallet.bat")
	} else {
		executable = filepath.Join(subDir, "boxwallet", "bin", "boxwallet")
	}

	// 2. Setup Context to ensure Elixir dies if Go is killed
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 3. Prepare the "start" command (Releases use 'start' or 'foreground')
	// 'foreground' is often better for Go to capture logs directly.
	cmd := exec.CommandContext(ctx, executable, "start")

	// 4. Pipe logs to Go's terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Starting Elixir engine on %s...\n", runtime.GOOS)

	if err := cmd.Run(); err != nil {
		fmt.Printf("Engine stopped: %v\n", err)
	}
}

func createScript(subDir string) {
	var scriptName, scriptContent string

	if runtime.GOOS == "windows" {
		scriptName = "run_boxwallet.bat"
		scriptContent = "@echo off\n.\\boxwallet\\bin\\boxwallet.bat start\npause"
	} else {
		scriptName = "run_boxwallet.sh"
		scriptContent = "#!/bin/bash\n./boxwallet/bin/boxwallet start"
	}

	fullPath := path.Join(subDir, scriptName)

	// 1. Create the script file
	err := os.WriteFile(fullPath, []byte(scriptContent), 0755)
	if err != nil {
		fmt.Printf("Error creating launcher: %v\n", err)
		return
	}

	fmt.Printf("\n✅ Launcher created: %s\n", fullPath)
	fmt.Println("Please run this file to start the BoxWallet engine.")
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
		version: "0.0.2",
	}

	// Immediate log
	fmt.Printf("%s%s%s v%s%s%s starting...\n",
		green, app.name, reset,
		green, app.version, reset,
	)

	// Defer runs when the surrounding function exits
	defer fmt.Printf("\n%s%s%s finished!\n",
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
		input = strings.TrimSpace(input)
		if input != "y" {
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

		if err := os.RemoveAll(dest); err != nil {
			fmt.Fprintf(os.Stderr, "unable to remove file: %v - %v", dest, err)
		}

		fmt.Println("All done!")
		createScript(latest_version)
		// fmt.Println("\n./" + dir + "/boxwallet/bin/boxwallet start")
	} else {
		// The latest version has already been downloaded, as the directory exists
		fmt.Println("You already have the latest version of " + green + "BoxWallet" + reset)
		createScript(latest_version)
	}

	// fmt.Printf("Latest release found: %s%s%s\n", green, downloadUrl, reset)
}
