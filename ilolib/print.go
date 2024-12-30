package ilolib

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/term"
)

type Printer interface {
	Print(log *log.Logger)
}

type HorizontalRule struct {
	Prefix string
	Suffix string
}

type InfoBox [][]string

func getTermWidth() int {
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		width = 50
	}
	return width
}

func (hr HorizontalRule) Print(log *log.Logger) {
	width := getTermWidth() - len(log.Prefix())
	log.Println(hr.Prefix + strings.Repeat("─", width-len(hr.Prefix)-len(hr.Suffix)) + hr.Suffix)
}

func (b InfoBox) Print(log *log.Logger) {
	width := getTermWidth() - len(log.Prefix())
	log.Printf(`╔%s╗`, strings.Repeat("═", width-2))

	separator := fmt.Sprintf("╟%s╢\n", strings.Repeat("─", width-2))

	for i, lines := range b {
		for _, line := range lines {
			log.Printf("║ %-*s ║\n", width-4, line)
		}
		if i+1 < len(b) {
			log.Print(separator)
		}
	}

	log.Printf(`╚%s╝`, strings.Repeat("═", width-2))
}
