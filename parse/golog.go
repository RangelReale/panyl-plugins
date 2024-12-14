package parse

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/RangelReale/panyl/v2"
)

var _ panyl.PluginParse = (*GoLog)(nil)

const GoLog_Format = "go_log"

// GoLog parse Golang log lines format
type GoLog struct {
	SourceAsCategory bool
}

// example: "2022-03-10T19:53:21.434Z	INFO	datadog-go/tracer.go:35	Datadog Tracer v1.28.0 ERROR: lost 2 traces"

var goLogRe = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T[^\s]+)\s+(\w+)\s+(.*\.go:\d+)\s+(.*)$`)

func (m *GoLog) ExtractParse(ctx context.Context, lines panyl.ProcessLines, result *panyl.Process) (bool, error) {
	// Only single line is supported
	if len(lines) != 1 {
		return false, nil
	}

	matches := goLogRe.FindStringSubmatch(lines.Line())
	if matches == nil {
		return false, nil
	}

	err := result.MergeLinesData(lines)
	if err != nil {
		return false, err
	}
	result.Line = ""

	timestamp := matches[1]
	level := matches[2]
	source := matches[3]
	message := matches[4]

	result.Data["timestamp"] = timestamp
	result.Data["level"] = level
	result.Data["source"] = source
	result.Data["message"] = message

	result.Metadata[panyl.MetadataFormat] = GoLog_Format
	result.Metadata[panyl.MetadataMessage] = message

	if timestamp != "" {
		ts, err := time.Parse(time.RFC3339Nano, timestamp)
		if err == nil {
			result.Metadata[panyl.MetadataTimestamp] = ts
		}
	}
	if level == "DEBUG" {
		result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
	} else if level == "INFO" {
		result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
	} else if level == "WARN" {
		result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
	} else if level == "ERROR" {
		result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelERROR
	}
	if m.SourceAsCategory && source != "" {
		result.Metadata[panyl.MetadataCategory] = strings.Split(source, ".")[0]
	}

	return true, nil
}

func (m GoLog) IsPanylPlugin() {}
