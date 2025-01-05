package toolbox

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/hairyhenderson/go-which"
)

type Toolbox map[string]string

func (t Toolbox) FindAndAdd(name string) error {
	programName := name

	if runtime.GOOS == "windows" && filepath.Ext(name) == "" {
		programName = name + ".exe"
	}

	path := which.Which(programName)
	if path == "" {
		return fmt.Errorf("add tool '%s': could not find on PATH", name)
	}

	t[name] = path
	return nil
}
