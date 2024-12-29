package main

type ToolInfo struct {
	Name string
	Path string
}

type Toolbox struct {
	Tools map[string]ToolInfo
}

func GetToolbox() (*Toolbox, error) {
	return &Toolbox{
		Tools: map[string]ToolInfo{
			"bash": {"Git Bash", `D:\Programs\Git\usr\bin\bash.exe`},
			"go":   {"Go", `C:\Program Files\Go\bin\go.exe`},
		},
	}, nil
}
