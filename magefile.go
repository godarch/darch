// +build mage

package main

import (
	"fmt"
	"os"
	"path"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	packageName = "github.com/pauldotknopf/darch/cmd/darch"
	ldflags     = "-X main.GitCommit=GITCOMMITTODO -X main.Version=VERSIONTODO"
)

// allow user to override go executable by running as GOEXE=xxx make ... on unix-like systems
var goexe = "go"

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}
}

// Install Go Dep and sync Hugo's vendored dependencies
func Vendor() error {
	return sh.Run("vndr")
}

func cleanBuild() error {
	fmt.Println("cleaning /bin")
	return removeDir("build")
}

// Build darch binary.
func Build() error {
	mg.Deps(cleanBuild)
	fmt.Println("building /bin/darch")
	return sh.Run(goexe, "build", "-ldflags", ldflags, "-o", "bin/darch", "pkg/cmd/darch/main.go")
}

func cleanBundle() error {
	fmt.Println("cleaning /bundle")
	return removeDir("bundle")
}

// Bundle the output into a format suitable for distribution (deb/rpm/etc).
func Bundle() error {
	mg.Deps(Build,
		cleanBundle)
	os.Mkdir("bundle", os.ModePerm)

	fmt.Println("placing bundle/bin/darch")
	if err := copyFile("bin/darch", "bundle/bin/darch"); err != nil {
		return err
	}
	return nil
}

func removeDir(dir string) error {
	_, err := sh.Exec(nil, os.Stdout, os.Stderr, "rm", "-rf", dir)
	return err
}

func copyFile(src, dest string) error {
	// make sure destination folder exists
	destParent := path.Dir(dest)
	if _, err := sh.Exec(nil, os.Stdout, os.Stderr, "mkdir", "-p", destParent); err != nil {
		return err
	}

	if _, err := sh.Exec(nil, os.Stdout, os.Stderr, "cp", src, dest); err != nil {
		return err
	}

	return nil
}
