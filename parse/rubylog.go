package parse

import (
	"context"
	"regexp"
	"time"

	"github.com/RangelReale/panyl"
)

var _ panyl.PluginParse = (*RubyLog)(nil)

const RubyLog_Format = "ruby_log"

// RubyLog parses Ruby log lines format
type RubyLog struct {
}

// example: "I, [1999-03-03T02:34:24.895701 #19074]  INFO -- Main: info."

var (
	rubyLogRe           = regexp.MustCompile(`(\w), \[(\d{4}-\d{2}-\d{2}T[^\s]+)\s+(#\d+)]\s+(\w+)\s+--\s+([^:]*):\s+(.*)`)
	rubyTimestampFormat = "2006-01-02T15:04:05.999999999"
)

func (m *RubyLog) ExtractParse(ctx context.Context, lines panyl.ProcessLines, result *panyl.Process) (bool, error) {
	// Only single line is supported
	if len(lines) != 1 {
		return false, nil
	}

	matches := rubyLogRe.FindStringSubmatch(lines.Line())
	if matches == nil {
		return false, nil
	}

	err := result.MergeLinesData(lines)
	if err != nil {
		return false, err
	}
	result.Line = ""

	timestamp := matches[2]
	level := matches[4]
	message := matches[6]

	result.Data["severity_id"] = matches[1]
	result.Data["timestamp"] = timestamp
	result.Data["pid"] = matches[3]
	result.Data["severity"] = level
	result.Data["prog_name"] = matches[5]
	result.Data["message"] = message

	result.Metadata[panyl.MetadataFormat] = RubyLog_Format
	result.Metadata[panyl.MetadataMessage] = message

	if timestamp != "" {
		ts, err := time.Parse(rubyTimestampFormat, timestamp)
		if err == nil {
			result.Metadata[panyl.MetadataTimestamp] = ts
		}
	}

	// https://ruby-doc.org/stdlib-2.6.4/libdoc/logger/rdoc/Logger.html
	if level == "DEBUG" {
		result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
	} else if level == "INFO" {
		result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
	} else if level == "WARN" {
		result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
	} else if level == "ERROR" || level == "FATAL" {
		result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelERROR
	}

	return true, nil
}

func (m RubyLog) IsPanylPlugin() {}
