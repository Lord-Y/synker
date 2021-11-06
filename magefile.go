//go:build mage

package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/magefile/mage/mg"
)

var (
	packageName  = "github.com/Lord-Y/synker"
	app          = "synker"
	artifacts    = "artifacts"
	sha          = "sha256sums.txt"
	architecture = map[string]string{
		"linux":   "amd64",
		"darwin":  "amd64",
		"windows": "amd64",
	}
)

var ldflags = fmt.Sprintf(
	"-X '%s/cmd.Version=%s' -X '%s/cmd.revision=%s' -X '%s/cmd.buildDate=%s' -X '%s/cmd.goVersion=%s'",
	packageName,
	os.Getenv("BUILD_VERSION"),
	packageName,
	os.Getenv("BUILD_REVISION"),
	packageName,
	time.Now().Format("2006-01-02T15:04:05Z0700"),
	packageName,
	runtime.Version(),
)

// Build Build binaries depending on os/architectures
func Build() (err error) {
	mg.Deps(InstallDeps)

	wg := sync.WaitGroup{}
	wg.Add(len(architecture))
	for oses, arch := range architecture {
		go func(oses, arch string, wg *sync.WaitGroup) (err error) {
			defer wg.Done()
			var appName string
			if oses == "windows" {
				appName = fmt.Sprintf("%s_%s_%s.exe", app, oses, arch)
			} else {
				appName = fmt.Sprintf("%s_%s_%s", app, oses, arch)
			}
			os.Remove(appName)
			fmt.Printf("Building %s ...\n", appName)
			cmd := exec.Command(
				"go",
				"build",
				"-ldflags",
				fmt.Sprintf(`%s`, ldflags),
				"-o",
				appName,
				".",
			)
			cmd.Env = append(
				os.Environ(),
				fmt.Sprintf("GOOS=%s", oses),
				fmt.Sprintf("GOARCH=%s", arch),
			)
			return cmd.Run()
		}(oses, arch, &wg)
	}
	wg.Wait()
	return
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	os.Setenv("GO111MODULE", "on")
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "mod", "download")
	return cmd.Run()
}

// Clean Remove all generated binaries
func Clean() {
	wg := sync.WaitGroup{}
	wg.Add(len(architecture))
	for oses, arch := range architecture {
		go func(oses, arch string, wg *sync.WaitGroup) {
			defer wg.Done()
			var appName string
			if oses == "windows" {
				appName = fmt.Sprintf("%s_%s_%s.exe", app, oses, arch)
			} else {
				appName = fmt.Sprintf("%s_%s_%s", app, oses, arch)
			}
			fmt.Printf("Cleaning %s ...\n", appName)
			os.Remove(appName)
			os.RemoveAll(artifacts)
		}(oses, arch, &wg)
	}
	wg.Wait()
}

// Compress Compress will gzip all binaries and generate checksum file
func Compress() (err error) {
	var (
		f, fsha *os.File
		writer  *gzip.Writer
		body    []byte
		appName string
	)
	_ = os.RemoveAll(artifacts)
	err = os.MkdirAll("artifacts", 0755)
	if err != nil {
		return err
	}

	for oses, arch := range architecture {
		if oses == "windows" {
			appName = fmt.Sprintf("%s_%s_%s.exe", app, oses, arch)
		} else {
			appName = fmt.Sprintf("%s_%s_%s", app, oses, arch)
		}

		f, err = os.OpenFile(fmt.Sprintf("%s/%s.tar.gz", artifacts, appName), os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		writer, err = gzip.NewWriterLevel(f, gzip.BestCompression)
		if err != nil {
			return err
		}
		defer writer.Close()

		tf := tar.NewWriter(writer)
		defer tf.Close()

		body, err = os.ReadFile(appName)
		if err != nil {
			return err
		}

		if body != nil {
			header := &tar.Header{
				Name: appName,
				Mode: int64(0644),
				Size: int64(len(body)),
			}
			err = tf.WriteHeader(header)
			if err != nil {
				return err
			}
			_, err := tf.Write(body)
			if err != nil {
				return err
			}
		}

		fsha, err = os.OpenFile(fmt.Sprintf("%s/%s", artifacts, sha), os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer fsha.Close()

		sum := sha256.Sum256(body)
		fsha.WriteString(fmt.Sprintf("%x\t%s\n", sum, appName))
	}
	return
}
