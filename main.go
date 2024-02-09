package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Check if gosec is installed
	if !isGosecInstalled() {
		fmt.Println("gosec is not installed. Attempting to install...")
		if err := installGosec(); err != nil {
			fmt.Println("Error installing gosec:", err)
			os.Exit(1)
		}

		// Retry to find gosec in PATH
		if !isGosecInstalled() {
			fmt.Println("gosec is not found in PATH even after installation.")
			os.Exit(1)
		}
	}

	// Read directory path from command-line arguments
	if len(os.Args) != 2 {
		fmt.Println("Usage: gosec-tool <directory-path>")
		os.Exit(1)
	}
	dirPath := os.Args[1]

	// Run gosec command to analyze the directory
	cmd := exec.Command("gosec", "-exclude=G104", "-fmt=json", "-out=gosec-report.json", dirPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		os.Exit(1)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting gosec command:", err)
		os.Exit(1)
	}

	// Read gosec output line by line and print it to the console
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for gosec command to finish:", err)
		os.Exit(1)
	}

	// Check if there were any errors reported by gosec
	if cmd.ProcessState.ExitCode() != 0 {
		fmt.Println("gosec command reported errors. See gosec-report.json for details.")
		os.Exit(1)
	}

	fmt.Println("Security analysis completed successfully. No issues found.")

}

func isGosecInstalled() bool {
	// Check if gosec is installed and available in PATH
	cmd := exec.Command("gosec", "--version")
	err := cmd.Run()
	return err == nil
}

func installGosec() error {
	// Install gosec using go get
	cmd := exec.Command("go", "install", "github.com/securego/gosec/v2/cmd/gosec")
	err := cmd.Run()
	if err != nil {
		return err
	}

	// Attempt to refresh PATH after installation
	if err := os.Setenv("PATH", os.Getenv("PATH")+":"+os.Getenv("GOPATH")+"/bin"); err != nil {
		return err
	}

	return nil
}
