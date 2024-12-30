package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"fourls.dev/ilo/ilolib"
)

func printHeader(name string, path string) {
	var log = log.New(os.Stdout, "", 0)
	log.Printf(`╔%s╗`, strings.Repeat("═", 50))
	log.Printf("║ ▒▒ %-45s ║\n", name)
	log.Printf(`╟%s╢`, strings.Repeat("─", 50))
	log.Printf("║ %-48s ║\n", path)
	log.Printf(`╚%s╝`, strings.Repeat("═", 50))
}

func main() {
	var wd, _ = os.Getwd()

	var path = flag.String("file", wd, "path to project definition file")
	var chosenFlow = flag.String("flow", "*", "name of flow to run")
	flag.Parse()

	var stat, _ = os.Stat(*path)
	if stat.IsDir() {
		*path = filepath.Join(*path, "ilo.yml")
	}

	if !filepath.IsAbs(*path) {
		var err error
		*path, err = filepath.Abs(*path)
		if err != nil {
			panic(err)
		}
	}

	file_contents, err := os.ReadFile(*path)
	if err != nil {
		panic(err)
	}

	project, err := ilolib.ParseYamlProjDef(file_contents)
	if err != nil {
		panic(err)
	}

	printHeader(project.Name, *path)

	toolbox, err := ilolib.GetToolbox()
	if err != nil {
		println("Error parsing tools.json: " + err.Error())
	}

	for _, flow := range project.Flows {
		if *chosenFlow == "*" || flow.Name == *chosenFlow {
			ilolib.ExecuteFlow(flow, *toolbox)
		}
	}
}
