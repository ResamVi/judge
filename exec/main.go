package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed submission
var sourceCode string

func main() {

	// Create a temporary directory for the build
	tmpDir, err := os.MkdirTemp("", "gorunner")
	if err != nil {
		fmt.Printf("Error creating temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir) // clean up temp files

	fmt.Println(tmpDir)

	source := filepath.Join(tmpDir, "main.go")
	binary := filepath.Join(tmpDir, "binary")

	err = os.WriteFile(source, []byte(sourceCode), 0644)
	if err != nil {
		panic(err)
	}

	buildCmd := exec.Command("go", "build", "-o", binary, source)

	var buildStderr bytes.Buffer
	buildCmd.Stderr = &buildStderr

	if err := buildCmd.Run(); err != nil {
		fmt.Printf("Build failed:\n%s\n", buildStderr.String())
		os.Exit(1)
	}

	// Run the compiled binary
	runCmd := exec.Command(binary)
	var stdout, stderr bytes.Buffer
	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr

	if err := runCmd.Run(); err != nil {
		fmt.Printf("Program failed: %v\nStderr:\n%s\n", err, stderr.String())
		os.Exit(1)
	}

	// Print the binary output
	fmt.Print(stdout.String())
}
