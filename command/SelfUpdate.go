package command

import (
	"fmt"
	"context"
	"log"
	"runtime"
	"strings"
	"github.com/google/go-github/github"
	"github.com/inconshreveable/go-update"
	"net/http"
)

type SelfUpdate struct {
	GithubOrganization  string
	GithubRepository    string
	GithubAssetTemplate string
}

func (conf *SelfUpdate) Execute(args []string) error {
	fmt.Println("Starting self update")

	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), conf.GithubOrganization, conf.GithubRepository)

	if _, ok := err.(*github.RateLimitError); ok {
		log.Println("GitHub rate limit, please try again later")
	}

	os := runtime.GOOS
	switch (runtime.GOOS) {
	case "darwin":
		os = "osx"
	}

	arch := runtime.GOARCH
	switch (arch) {
	case "amd64":
		arch = "x64"
	case "386":
		arch = "x32"
	}
	assetName := conf.GithubAssetTemplate
	assetName = strings.Replace(assetName, "%OS%", os, -1)
	assetName = strings.Replace(assetName, "%ARCH%", arch, -1)

	fmt.Println(fmt.Sprintf(" - searching for asset \"%s\"", assetName))

	for _, asset := range release.Assets {
		if asset.GetName() == assetName {
			downloadUrl := asset.GetBrowserDownloadURL()
			fmt.Println(fmt.Sprintf(" - found new update url \"%s\"", downloadUrl))
			conf.runUpdate(downloadUrl)
			fmt.Println(fmt.Sprintf(" - finished update to version %s", release.GetName()))
			return nil
		}
	}

	fmt.Println(" - unable to download latest version")
	return nil
}

func (conf *SelfUpdate) runUpdate(url string) error {
	fmt.Println(" - downloading update")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println(" - applying update")
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		// error handling
	}
	return err
}
