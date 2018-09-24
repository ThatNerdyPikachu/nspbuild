package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	VERSION = "1.0"
)

type release struct {
	Assets []asset `json:"assets"`
}

type asset struct {
	URL string `json:"browser_download_url"`
}

type nacp struct {
	Name    string `json:"name"`
	Author  string `json:"author"`
	Version string `json:"version"`
	TitleID string `json:"title_id"`
}

func printHelpAndExit() {
	lines := []string{
		fmt.Sprintf("nspbuild v%s by Pika", VERSION),
		"usage: nspbuild <path/to/nso> <name> <author> <version> <path/to/icon/jpg> <tid>",
	}

	for _, v := range lines {
		fmt.Printf("%s\n", v)
	}

	os.Exit(0)
}

func getRelease(repo string) (release, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo))
	if err != nil {
		return release{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return release{}, err
	}
	resp.Body.Close()

	r := release{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return release{}, err
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
		"version",
		"icon",
		"tid",
	})

	if len(args) != 6 {
		printHelpAndExit()
	}

	if !fileExists(args["nso"]) {
		fmt.Printf("error: the file at %s does not exist!\n", args["nso"])
		os.Exit(1)
	} else if args["icon"] != "none" && !fileExists(args["icon"]) {
		fmt.Printf("error: the file at %s does not exist!\n", args["icon"])
		os.Exit(1)
	}

	if !isHex(args["tid"]) {
		fmt.Printf("error: the title id %s is not valid hex!\n", args["tid"])
		os.Exit(1)
	} else if len(args["tid"]) != 16 {
		fmt.Printf("error: the title id %s is not 16 characters!\n", args["tid"])
		os.Exit(1)
	}

	chkErr(os.MkdirAll("build/", 0700))

	linkle, err := getRelease("MegatonHammer/linkle")
	chkErr(err)

	s := ""
	if runtime.GOOS == "windows" && runtime.GOARCH == "amd64" {
		s = "x86_64-pc-windows-msvc.zip"
	}

	if s == "" {
		fmt.Printf("error: nspbuild does not support this os/arch as of now, exiting...\n")
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

	chkErr(download(url, "build/linkle.zip"))

	chkErr(unzipFile("build/linkle.zip", "linkle.exe", "build/linkle.exe"))

	hbp, err := getRelease("The-4n/hacBrewPack")
	chkErr(err)

	chkErr(download(hbp.Assets[0].URL, "build/hbp.zip"))

	chkErr(unzipFile("build/hbp.zip", "hacbrewpack.exe", "build/hbp.exe"))

	chkErr(download("https://raw.githubusercontent.com/ThatNerdyPikachu/nspbuild/master/npdmtool.exe",
		"build/npdmtool.exe"))

	chkErr(os.MkdirAll("build/exefs", 0700))

	chkErr(copy(args["nso"], "build/exefs/main"))

	resp, err := http.Get("https://raw.githubusercontent.com/switchbrew/nx-hbloader/master/hbl.json")
	chkErr(err)

	npdm, err := os.Create("build/npdm.json")
	chkErr(err)

	scanner := bufio.NewScanner(resp.Body)

	replacer := strings.NewReplacer("hbloader", args["name"], "0x010000000000100D", "0x"+strings.ToLower(args["tid"]))

	for scanner.Scan() {
		npdm.WriteString(replacer.Replace(scanner.Text()) + "\n")
	}

	resp.Body.Close()
	npdm.Close()

	cmd := exec.Command(".\\npdmtool", "npdm.json", "exefs/main.npdm")
	cmd.Dir = "build/"
	chkErr(cmd.Run())

	chkErr(os.MkdirAll("build/control", 0700))

	gen := nacp{
		Name:    args["name"],
		Author:  args["author"],
		Version: args["version"],
		TitleID: strings.ToLower(args["tid"]),
	}

	nacp, err := os.Create("build/nacp.json")
	chkErr(err)

	j, err := json.Marshal(gen)
	chkErr(err)

	nacp.WriteString(string(j))
	nacp.Close()

	cmd = exec.Command(".\\linkle", "nacp", "nacp.json", "control/control.nacp")
	cmd.Dir = "build/"
	chkErr(cmd.Run())

	if args["icon"] != "none" {
		languages := []string{
			"Japanese",
			"AmericanEnglish",
			"French",
			"German",
			"Italian",
			"Spanish",
			"Chinese",
			"Korean",
			"Dutch",
			"Portuguese",
			"Russian",
			"Taiwanese",
			"BritishEnglish",
			"CanadianFrench",
			"LatinAmericanSpanish",
			"SimplifiedChinese",
			"TraditionalChinese",
		}

		for _, v := range languages {
			chkErr(copy(args["icon"], fmt.Sprintf("build/control/icon_%s.dat", v)))
		}
	}

	chkErr(copy("keys.txt", "build/keys.dat"))

	cmd = exec.Command(".\\hbp", "--noromfs", "--nologo")
	cmd.Dir = "build/"
	chkErr(cmd.Run())

	chkErr(os.MkdirAll("out/", 0700))

	chkErr(copy(fmt.Sprintf("build/hacbrewpack_nsp/%s.nsp", strings.ToLower(args["tid"])),
		fmt.Sprintf("out/%s [%s].nsp", args["name"], strings.ToLower(args["tid"]))))

	chkErr(os.RemoveAll("build/"))
}
