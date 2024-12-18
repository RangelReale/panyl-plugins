package parse

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/RangelReale/panyl/v2"
)

const GoSpaceLogFormat = "go_space_log"

// GoSpaceLog parse Golang log lines format
type GoSpaceLog struct {
	SourceAsCategory bool
}

var _ panyl.PluginParse = GoSpaceLog{}

// example: "level=info ts=2024-12-18T14:55:27.787558447Z caller=poller.go:136 msg="blocklist poll complete" seconds=0.128523552"

var GoSpaceLogRe = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T[^\s]+)\s+(\w+)\s+(.*\.go:\d+)\s+(.*)$`)

func (m GoSpaceLog) ExtractParse(ctx context.Context, lines panyl.ItemLines, item *panyl.Item) (bool, error) {
	// Only single line is supported
	if len(lines) != 1 {
		return false, nil
	}

	fields := m.splitFields(lines.Line())
	if len(fields) == 0 {
		return false, nil
	}

	for _, fieldToCheck := range []string{"level", "ts", "msg"} {
		if _, ok := fields[fieldToCheck]; !ok {
			return false, nil
		}
	}

	err := item.MergeLinesData(lines)
	if err != nil {
		return false, err
	}

	for fn, fv := range fields {
		item.Data[fn] = fv
	}

	item.Line = ""

	item.Metadata[panyl.MetadataFormat] = GoSpaceLogFormat

	if timestamp, ok := fields["ts"]; ok {
		ts, err := time.Parse(time.RFC3339Nano, timestamp)
		if err == nil {
			item.Metadata[panyl.MetadataTimestamp] = ts
		}
	}
	if level, ok := fields["level"]; ok {
		if level == "debug" {
			item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
		} else if level == "info" {
			item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
		} else if level == "warn" || level == "warning" {
			item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
		} else if level == "error" {
			item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelERROR
		}
	}
	if message, ok := fields["msg"]; ok {
		item.Metadata[panyl.MetadataMessage] = message
	}
	if source, ok := fields["caller"]; ok {
		if m.SourceAsCategory && source != "" {
			item.Metadata[panyl.MetadataCategory] = strings.Split(source, ":")[0]
		}
	}

	return true, nil
}

func (m GoSpaceLog) IsPanylPlugin() {}

func (m GoSpaceLog) splitFields(str string) map[string]string {
	fields := m.splitString(str)
	if len(fields) == 0 {
		return nil
	}
	ret := make(map[string]string)
	for _, field := range fields {
		name, value, ok := strings.Cut(field, "=")
		if ok {
			ret[name] = m.trimQuotes(value)
		}
	}
	return ret
}

func (m GoSpaceLog) trimQuotes(str string) string {
	s, err := strconv.Unquote(strings.TrimSpace(str))
	if err != nil {
		return str
	}
	return s
}

func (m GoSpaceLog) splitString(str string) []string {
	quoted := false
	backTick := false
	return strings.FieldsFunc(str, func(r1 rune) bool {
		if backTick {
			backTick = false
		} else if r1 == '\\' {
			backTick = true
		} else if r1 == '"' {
			quoted = !quoted
		}
		return !quoted && unicode.IsSpace(r1)
	})
}
