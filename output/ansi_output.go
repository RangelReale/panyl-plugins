package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/RangelReale/panyl"
	"github.com/fatih/color"
	"time"
)

type AnsiOutputSprintfFunc func(format string, a ...interface{}) string

type AnsiOutput struct {
	ColorInformation, ColorWarning, ColorError, ColorInternalError, ColorUnknown AnsiOutputSprintfFunc
}

func NewAnsiOutput(ansi bool) *AnsiOutput {
	ret := &AnsiOutput{
		ColorError:         fmt.Sprintf,
		ColorWarning:       fmt.Sprintf,
		ColorInformation:   fmt.Sprintf,
		ColorInternalError: fmt.Sprintf,
		ColorUnknown:       fmt.Sprintf,
	}
	if ansi {
		ret.ColorError = color.New(color.FgRed).SprintfFunc()
		ret.ColorWarning = color.New(color.FgYellow).SprintfFunc()
		ret.ColorInformation = color.New(color.FgGreen).SprintfFunc()
		ret.ColorInternalError = color.New(color.FgHiRed).SprintfFunc()
		ret.ColorUnknown = color.New(color.FgMagenta).SprintfFunc()
	}
	return ret
}

func (o *AnsiOutput) OnResult(p *panyl.Process) (cont bool) {
	var out bytes.Buffer

	// level
	var levelColor AnsiOutputSprintfFunc
	level := p.Metadata.StringValue(panyl.Metadata_Level)
	switch level {
	case panyl.MetadataLevel_DEBUG, panyl.MetadataLevel_INFO:
		levelColor = o.ColorInformation
	case panyl.MetadataLevel_WARNING:
		levelColor = o.ColorWarning
	case panyl.MetadataLevel_CRITICAL, panyl.MetadataLevel_FATAL:
		levelColor = o.ColorError
	default:
		level = "unknown"
		levelColor = o.ColorUnknown
	}

	// timestamp
	if ts, ok := p.Metadata[panyl.Metadata_Timestamp]; ok {
		out.WriteString(fmt.Sprintf("%s ", ts.(time.Time).Local().Format("2006-01-02 15:04:05.000")))
	}

	// application
	if application := p.Metadata.StringValue(panyl.Metadata_Application); application != "" {
		out.WriteString(fmt.Sprintf("| %s | ", application))
	}

	// level
	if level != "" {
		out.WriteString(fmt.Sprintf("[%s] ", level))
	}

	// format
	if format := p.Metadata.StringValue(panyl.Metadata_Format); format != "" {
		out.WriteString(fmt.Sprintf("(%s) ", format))
	}

	// category
	if category := p.Metadata.StringValue(panyl.Metadata_Category); category != "" {
		out.WriteString(fmt.Sprintf("{{%s}} ", category))
	}

	// message
	if msg := p.Metadata.StringValue(panyl.Metadata_Message); msg != "" {
		out.WriteString(msg)
	} else if p.Line != "" {
		out.WriteString(p.Line)
	} else if len(p.Data) > 0 {
		dt, err := json.Marshal(p.Data)
		if err != nil {
			fmt.Println(o.ColorInternalError("Error marshaling data to json: %s", err.Error()))
			return
		}
		out.WriteString(fmt.Sprintf("| %s", string(dt)))
	}

	fmt.Println(levelColor(out.String()))

	return true
}
