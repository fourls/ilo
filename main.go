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

	for i, flow := range data.Flows {
		println("Running flow", i)
		RunFlow(flow)
	}
}
