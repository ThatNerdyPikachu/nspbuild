package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const (
	VERSION = "0.1"
)

type Release struct {
	Assets []Asset `json:"assets"`
}

type Asset struct {
	URL string `json:"browser_download_url"`
}

func printHelpAndExit() {
	lines := []string{
		fmt.Sprintf("nspbuild v%s by Pika", VERSION),
		"usage: nspbuild <path/to/nso> <name> <author> <path/to/icon/jpg> <tid>",
	}

	for _, v := range lines {
		fmt.Printf("%s\n", v)
	}

	os.Exit(0)
}

func getRelease(repo string) (Release, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo))
	if err != nil {
		return Release{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Release{}, err
	}
	resp.Body.Close()

	r := Release{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return Release{}, err
	}

	return r, nil
}

func main() {
	os.Args = del(os.Args, 0)
	if isEverythingNil(os.Args) {
		printHelpAndExit()
	}

	args := org(parse(os.Args), []string{
		"nso",
		"name",
		"author",
		"icon",
		"tid",
	})

	if !fileExists(args["nso"]) {
		fmt.Printf("error: the file at %s does not exist!\n", args["nso"])
		os.Exit(1)
	} /* else if !fileExists(args["icon"]) {
		fmt.Printf("error: the file at %s does not exist!", args["icon"])
		os.Exit(1)
	} */

	if !isHex(args["tid"]) {
		fmt.Printf("error: the title id %s is not valid hex!\n", args["tid"])
		os.Exit(1)
	} else if len(args["tid"]) != 16 {
		fmt.Printf("error: the title id %s is not 16 characters!\n", args["tid"])
	}

	err := os.MkdirAll("build/", 0700)
	if err != nil {
		panic(err)
	}

	linkle, err := getRelease("MegatonHammer/linkle")
	if err != nil {
		panic(err)
	}

	s := ""
	if runtime.GOOS == "windows" && runtime.GOARCH == "386" {
		// TODO: verify this arch
		s = "i686-pc-windows-msvc.zip"
	} else if runtime.GOOS == "windows" && runtime.GOARCH == "amd64" {
		s = "x86_64-pc-windows-msvc.zip"
	}

	if s == "" {
		fmt.Printf("error: linkle does not support this os/arch as of now, exiting...\n")
		os.Exit(1)
	}

	url := ""
	for _, v := range linkle.Assets {
		if strings.HasSuffix(v.URL, s) {
			url = v.URL
		}
	}

	if url == "" {
		fmt.Printf("error: could not find a linke build for your os/arch, exiting...\n")
		os.Exit(1)
	}

	err = download(url, "build/linkle.zip")
	if err != nil {
		panic(err)
	}

	err = unzip("build/linkle.zip", "linkle.exe", "build/linkle.exe")
	if err != nil {
		panic(err)
	}

	hbp, err := getRelease("The-4n/hacBrewPack")
	if err != nil {
		panic(err)
	}

	err = download(hbp.Assets[0].URL, "build/hbp.zip")
	if err != nil {
		panic(err)
	}

	err = unzip("build/hbp.zip", "hacbrewpack.exe", "build/hbp.exe")
	if err != nil {
		panic(err)
	}

	err = download("https://raw.githubusercontent.com/ThatNerdyPikachu/nspbuild/master/npdmtool.exe",
		"build/npdmtool.exe")
	if err != nil {
		panic(err)
	}

	if isAnythingNil(mapToSlice(args)) {
		printHelpAndExit()
	}
}
