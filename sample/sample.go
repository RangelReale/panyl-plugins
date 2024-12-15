package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/RangelReale/panyl-plugins/v2/metadata"
	"github.com/RangelReale/panyl-plugins/v2/parse"
	"github.com/RangelReale/panyl/v2"
	"github.com/RangelReale/panyl/v2/plugins/clean"
	"github.com/RangelReale/panyl/v2/plugins/structure"
)

func main() {
	ctx := context.Background()

	processor := panyl.NewProcessor(
		panyl.WithPlugins(
			&clean.AnsiEscape{},
			&metadata.DockerCompose{},
			&structure.JSON{},
			&parse.GoLog{},
			&parse.RubyLog{},
			&parse.MongoLog{},
			&parse.NGINXErrorLog{},
		),
		// may use a logger when debugging, it outputs each source line and parsed processes
		// panyl.WithLogger(panyl.NewStdLogOutput()),
	)

	err := processor.Process(ctx, os.Stdin, &Output{}, panyl.WithLineLimit(0, 100))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing input: %s", err.Error())
	}
}

type Output struct {
}

func (o *Output) OnResult(ctx context.Context, p *panyl.Process) (cont bool) {
	var out bytes.Buffer

	// timestamp
	if ts, ok := p.Metadata[panyl.MetadataTimestamp]; ok {
		out.WriteString(fmt.Sprintf("%s ", ts.(time.Time).Local().Format("2006-01-02 15:04:05.000")))
	}

	// level
	if level := p.Metadata.StringValue(panyl.MetadataLevel); level != "" {
		out.WriteString(fmt.Sprintf("[%s] ", level))
	}

	// category
	if category := p.Metadata.StringValue(panyl.MetadataCategory); category != "" {
		out.WriteString(fmt.Sprintf("{{%s}} ", category))
	}

	// message
	if msg := p.Metadata.StringValue(panyl.MetadataMessage); msg != "" {
		out.WriteString(msg)
	} else if len(p.Data) > 0 {
		// Extracted structure but no metadata
		dt, err := json.Marshal(p.Data)
		if err != nil {
			fmt.Printf("Error marshaling data to json: %s\n", err.Error())
			return
		}
		out.WriteString(fmt.Sprintf("| %s", string(dt)))
	} else if p.Line != "" {
		// Show raw line if available
		out.WriteString(p.Line)
	}

	fmt.Println(out.String())
	return true
}

func (o *Output) OnFlush(ctx context.Context) {}

func (o *Output) OnClose(ctx context.Context) {}
