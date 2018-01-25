// +build mage

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	packageName = "github.com/pauldotknopf/darch/cmd/darch"
	ldflags     = "-X main.GitCommit=$HASH -X main.Version=$VERSION"
)

// allow user to override go executable by running as GOEXE=xxx make ... on unix-like systems
var goexe = "go"
var goarch = "amd64"
var version = ""
var hash = ""
var tag = ""
var isTagBuild = false

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// Update the GOARCH.
	if output, err := sh.Output(goexe, "env", "GOARCH"); err == nil {
		goarch = output
	} else {
		panic(fmt.Sprintf("couldn't get GOARCH: %v", err))
	}

	// Get the version number for the build.
	travisTag, exists := os.LookupEnv("TRAVIS_TAG")
	if exists && len(travisTag) > 0 {
		tag = travisTag
		isTagBuild = true
		if strings.HasPrefix(travisTag, "v") {
			version = travisTag[1:]
		} else {
			version = travisTag
		}
	} else {
		// TODO: get current tag, if any. Or, use gitversion.
		version = "NA"
	}

	if output, err := sh.Output("git", "rev-parse", "--short", "HEAD"); err == nil {
		hash = output
	} else {
		hash = "NA"
	}

	fmt.Printf("version: %v\n", version)
	fmt.Printf("hash: %v\n", hash)
	fmt.Printf("tag: %v\n", tag)
	fmt.Printf("istagbuild: %v\n", isTagBuild)
}

// Install Go Dep and sync Hugo's vendored dependencies
func Vendor() error {
	mg.Deps(ensureVndrTool)
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
	return sh.RunWith(map[string]string{
		"HASH":    hash,
		"VERSION": version,
	}, goexe, "build", "-ldflags", ldflags, "-o", "bin/darch", "pkg/cmd/darch/main.go")
}

func Test() error {
	return sh.Run(goexe, "test", "./pkg/...")
}

func Lint() error {
	mg.Deps(ensureGoMetaLinterTool)
	return sh.Run("gometalinter.v2", "--config", ".gometalinter.json", "./pkg/...")
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
	if err := copyFile("bin/darch", "bundle/usr/bin/darch"); err != nil {
		return err
	}

	fmt.Println("placing bundle/etc/darch/hooks/fstab/hook")
	if err := copyFile("scripts/hooks/fstab", "bundle/etc/darch/hooks/fstab/hook"); err != nil {
		return err
	}

	fmt.Println("placing bundle/etc/darch/hooks/hostname/hook")
	if err := copyFile("scripts/hooks/hostname", "bundle/etc/darch/hooks/hostname/hook"); err != nil {
		return err
	}

	fmt.Println("placing bundle/etc/grub.d/60_darch")
	if err := copyFile("scripts/grub-mkconfig-script", "bundle/etc/grub.d/60_darch"); err != nil {
		return err
	}

	fmt.Println("packaging bundle")
	if _, err := sh.Exec(nil,
		os.Stdout,
		os.Stderr,
		"tar",
		"cvzpf",
		fmt.Sprintf("bundle/darch-%s.tar.gz", goarch),
		"-C",
		"bundle",
		"usr/bin/darch",
		"etc/grub.d/60_darch",
		"etc/darch/hooks"); err != nil {
		return err
	}

	return nil
}

func Release() error {
	mg.Deps(
		ensureGithubRelease,
	)

	if !isTagBuild {
		fmt.Println("no a tag build, skipping release")
		return nil
	}

	fmt.Println("creating release in github")
	err := sh.Run("github-release",
		"release",
		"--user",
		"pauldotknopf",
		"--repo",
		"darch",
		"--tag",
		tag)

	if err != nil {
		return err
	}

	fmt.Println("uploading release to github")
	err = sh.Run("github-release",
		"upload",
		"--user",
		"pauldotknopf",
		"--repo",
		"darch",
		"--tag",
		tag,
		"--name",
		fmt.Sprintf("darch-%s.tar.gz", goarch),
		"--file",
		fmt.Sprintf("bundle/darch-%s.tar.gz", goarch))
	if err != nil {
		return err
	}

	fmt.Println("updating arch aur")
	return sh.Run("scripts/aur/deploy-aur", version)
}

func CI() error {
	mg.SerialDeps(
		Vendor,
		Lint,
		Bundle,
		Release,
	)
	return nil
}

var Default = Build
