// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"os"
	"os/exec"
	"path"

	"github.com/magefile/mage/sh"
)

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

	if _, err := sh.Exec(nil, os.Stdout, os.Stderr, "cp", "-p", src, dest); err != nil {
		return err
	}

	return nil
}

func ensureGoTool(tool, pkg string) error {
	_, err := exec.LookPath(tool)
	if err != nil {
		fmt.Printf("couldn't find tool %s, attempting to download\n", tool)
		if _, err := sh.Exec(nil, os.Stdout, os.Stderr, "go", "get", "-u", pkg); err != nil {
			return err
		}
		_, err := exec.LookPath(tool)
		if err != nil {
			return fmt.Errorf("couldn't download tool %s from %s", tool, pkg)
		}
	}
	return nil
}

func ensureGoLintTool() error {
	return ensureGoTool("golint", "github.com/golang/lint/golint")
}

func ensureGoMetaLinterTool() error {
	mg.Deps(
		ensureGoLintTool,
	)
	return ensureGoTool("gometalinter.v2", "gopkg.in/alecthomas/gometalinter.v2")
}

func ensureVndrTool() error {
	return ensureGoTool("vndr", "github.com/LK4D4/vndr")
}

func ensureGithubRelease() error {
	return ensureGoTool("github-release", "github.com/aktau/github-release")
}
