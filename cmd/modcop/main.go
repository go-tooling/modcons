package main

import (
	"flag"
	"fmt"
	. "github.com/go-tooling/modcop"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	rulesPath := flag.String("rulepath", "", "Path to rules file, may be local or http/s")
	modPath := flag.String("modpath", "./go.mod", "Path to rules file, may be local or http/s, defaults to ./go.mod")
	parseOnly := flag.Bool("parseOnly", false, "Only parse the rule file")

	flag.Parse()

	if len(*rulesPath) < 1 {
		showUsage()
		return
	}

	rulesRdr, err := getFile(*rulesPath)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer rulesRdr.Close()

	modRdr, err := getFile(*modPath)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer modRdr.Close()

	rules, err := ParseRules(rulesRdr)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	mods, err := ParseModFile(modRdr)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	if !*parseOnly {
		includesDeprecated := false

		for _, mod := range mods {

			if rules.IsDeprecated(mod) {
				includesDeprecated = true
				fmt.Printf("Deprecated mod: %v %v", mod.Path, mod.Version.String())
			} else {
				fmt.Printf("%v %v OK", mod.Path, mod.Version.String())
			}
		}

		if includesDeprecated {
			os.Exit(1)
		}
	}
}

func getFile(path string) (io.ReadCloser, error) {
	if strings.Index(path, "http") > -1 {
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	} else {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

func showUsage() {
	fmt.Println("Switches (prefix with -)")
	fmt.Println("\t rulepath: Path to rules file, may be local or http/s")
	fmt.Println("\t modpath: ath to rules file, may be local or http/s, defaults to ./go.mod")
}
