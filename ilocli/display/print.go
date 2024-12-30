package display

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
	Header string
	Footer string
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

	headerStart := -1
	headerEnd := -1
	footerStart := -1
	footerEnd := -1

	headerEdge := func(i int) bool {
		return headerStart == i || headerEnd == i
	}

	footerEdge := func(i int) bool {
		return footerStart == i || footerEnd == i
	}

	if hr.Header != "" {
		headerStart = 0
		headerEnd = len(hr.Header) + 3
		log.Printf("┌%s┐\n", strings.Repeat("─", len(hr.Header)+2))
		log.Printf("│ %s │\n", hr.Header)
	}

	if hr.Footer != "" {
		footerStart = 0
		footerEnd = len(hr.Footer) + 3
	}

	var sb strings.Builder

	for i := range width {
		switch {
		case !headerEdge(i) && !footerEdge(i):
			sb.WriteRune('─')
		case headerEdge(i) && !footerEdge(i):
			sb.WriteRune('┴')
		case !headerEdge(i) && footerEdge(i):
			sb.WriteRune('┬')
		case headerEdge(i) && footerEdge(i):
			sb.WriteRune('┼')
		}
	}

	log.Print(sb.String())

	if hr.Footer != "" {
		log.Printf("│ %s │\n", hr.Footer)
		log.Printf("└%s┘\n", strings.Repeat("─", len(hr.Footer)+2))
	}
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
