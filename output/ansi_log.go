package output

import (
	"bytes"
	"fmt"
	"github.com/RangelReale/panyl"
	"github.com/RangelReale/panyl/util"
	"github.com/fatih/color"
)

type AnsiLog struct {
	ShowSource bool
}

func (l AnsiLog) LogSourceLine(n int, line, rawLine string) {
	red := color.New(color.FgRed)

	red.Printf("@@@ SOURCE LINE [%d]: '%s' @@@\n", n, line)
}

func (l AnsiLog) LogProcess(p *panyl.Process) {
	green := color.New(color.FgGreen)

	var lineno string
	if p.LineCount > 1 {
		lineno = fmt.Sprintf("[%d-%d]", p.LineNo, p.LineNo+p.LineCount-1)
	} else {
		lineno = fmt.Sprintf("[%d]", p.LineNo)
	}

	var buf bytes.Buffer

	if len(p.Metadata) > 0 {
		_, _ = buf.WriteString(fmt.Sprintf("Metadata: %+v", p.Metadata))
	}
	if len(p.Data) > 0 {
		if buf.Len() > 0 {
			_, _ = buf.WriteString(" - ")
		}
		_, _ = buf.WriteString(fmt.Sprintf("Data: %+v", p.Data))
	}

	if len(p.Line) > 0 {
		if buf.Len() > 0 {
			_, _ = buf.WriteString(" - ")
		}
		_, _ = buf.WriteString(fmt.Sprintf("Line: \"%s\"", p.Line))
	}

	if l.ShowSource && len(p.Source) > 0 {
		if buf.Len() > 0 {
			_, _ = buf.WriteString(" - ")
		}
		_, _ = buf.WriteString(fmt.Sprintf("Source: \"%s\"", util.DoAnsiEscapeString(p.Source)))
	}

	green.Printf("*** PROCESS LINE %s: %s\n", lineno, buf.String())
}
