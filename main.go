package main

import (
	"os"
)

func main() {
	file_contents, err := os.ReadFile("test.ilo.json")
	if err != nil {
		panic(err)
	}

	data, err := ParseProjDef(file_contents)
	if err != nil {
		panic(err)
	}

	for _, tool := range data.Tools {
		println("Requires", tool.Name, "version", tool.Version)
	}
}
