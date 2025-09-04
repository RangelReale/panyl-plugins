package parse

import (
	"context"
	"regexp"
	"time"

	"github.com/RangelReale/panyl/v2"
)

const ErlangLogFormat = "erlang_log"

// ErlangLog parses Erlang log lines format
type ErlangLog struct {
}

// example: "2025-09-04T13:26:52.315705+00:00 [warning] not_all_kafka_partitions_connected"

var _ panyl.PluginParse = ErlangLog{}

var (
	erlangLogRe           = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\S+)\s+\[(\w+)\]\s+(.*)$`)
	erlangTimestampFormat = "2006-01-02T15:04:05.999999-07:00"
)

func (m ErlangLog) ExtractParse(ctx context.Context, lines panyl.ItemLines, item *panyl.Item) (bool, error) {
	// Only single line is supported
	if len(lines) != 1 {
		return false, nil
	}

	matches := erlangLogRe.FindStringSubmatch(lines.Line())
	if matches == nil {
		return false, nil
	}

	err := item.MergeLinesData(lines)
	if err != nil {
		return false, err
	}
	item.Line = ""

	timestamp := matches[1]
	level := matches[2]
	message := matches[3]

	item.Data["timestamp"] = timestamp
	item.Data["level"] = level
	item.Data["message"] = message

	item.Metadata[panyl.MetadataFormat] = ErlangLogFormat
	item.Metadata[panyl.MetadataMessage] = message

	if timestamp != "" {
		ts, err := time.Parse(erlangTimestampFormat, timestamp)
		if err == nil {
			item.Metadata[panyl.MetadataTimestamp] = ts
		}
	}

	if level == "debug" {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
	} else if level == "warning" {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
	} else if level == "error" {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelERROR
	} else {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
	}

	return true, nil
}

func (m ErlangLog) IsPanylPlugin() {}
