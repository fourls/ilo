package main

import (
	"os"
)

func main() {
	file_contents, err := os.ReadFile("test.ilo.yaml")
	if err != nil {
		panic(err)
	}

	data, err := ParseYamlProjDef(file_contents)
	if err != nil {
		panic(err)
	}

	for _, flow := range data.Flows {
		ExecuteFlow(flow)
	}
}
