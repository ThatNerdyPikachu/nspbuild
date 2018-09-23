package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	VERSION = "0.1"
)

func printHelpAndExit() {
	lines := []string{
		fmt.Sprintf("nro2nso v%s by Pika", VERSION),
		"usage: nro2so <path/to/nro> <path/to/output/nso>",
	}

	for _, v := range lines {
		fmt.Printf("%s\n", v)
	}

	os.Exit(0)
}

func main() {
	os.Args = del(os.Args, 0)
	if isEverythingNil(os.Args) {
		printHelpAndExit()
	}

	args := org(parse(os.Args), []string{
		"nro",
		"nso",
	})

	if len(args) != 2 {
		printHelpAndExit()
	}

	if !fileExists(args["nro"]) {
		fmt.Printf("error: the file at %s does not exist!\n", args["nro"])
		os.Exit(1)
	}

	err := os.MkdirAll("temp/", 0700)
	if err != nil {
		panic(err)
	}

	err = download("https://raw.githubusercontent.com/ThatNerdyPikachu/nspbuild/master/binaries/nx2elf.exe",
		"temp/nxtool.exe")
	if err != nil {
		panic(err)
	}

	err = copy(args["nro"], "temp/app.nro")
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(".\\nxtool", "--elf=app.elf", "app.nro")
	cmd.Dir = "temp/"
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	err = download("https://raw.githubusercontent.com/ThatNerdyPikachu/nspbuild/master/binaries/nx2elf.exe",
		"temp/elf2nso.exe")
	if err != nil {
		panic(err)
	}

	cmd = exec.Command(".\\elf2nso", "app.elf", "app.nso")
	cmd.Dir = "temp/"
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	err = copy("temp/app.nso", args["nso"])
	if err != nil {
		panic(err)
	}

	err = os.RemoveAll("temp/")
	if err != nil {
		panic(err)
	}
}
