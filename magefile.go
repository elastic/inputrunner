//go:build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"os/exec"
)

// Format formats all source files with `go fmt`
func Format() error {
	return sh.RunV("go", "fmt", "./...")
}

// Build downloads dependencies and builds the inputrunner binary
func Build() error {
	if err := sh.RunV("go", "mod", "download"); err != nil {
		return err
	}

	return sh.RunV("go", "build", ".")
}

// Test runs all unit tests and writes an HTML coverage report to the build directory
func Test() error {
	err := sh.RunV("go", "test", "./...", "-coverprofile=coverage.out")
	sh.RunV("go", "tool", "cover", "-html=coverage.out", "-o", "build/coverage.html")
	return err
}

// Check runs static analysis and security checks
func Check() error {
	if installed := install("staticcheck", "honnef.co/go/tools/cmd/staticcheck@latest"); !installed {
		return nil
	}

	fmt.Println("Running staticcheck...")
	if err := sh.RunV("staticcheck", "-f=stylish", "./..."); err != nil {
		// fixme - lots of errors right now...
	}

	if installed := install("gosec", "github.com/securego/gosec/v2/cmd/gosec@latest"); !installed {
		return nil
	}

	fmt.Println("Running gosec...")
	if err := sh.RunV("gosec", "./..."); err != nil {
		// fixme - lots of errors right now...
	}

	return nil
}

func install(packageName, installURL string) (isInstalled bool) {
	_, missing := exec.LookPath(packageName)
	if missing != nil {
		fmt.Printf("installing %v...\n", packageName)
		err := sh.RunV("go", "install", installURL)
		if err != nil {
			fmt.Printf("Could not install %v, skipping...\n", packageName)
			return false
		}
		fmt.Printf("%v installed...\n", packageName)
	}
	return true
}