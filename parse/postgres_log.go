package parse

import (
	"context"
	"regexp"
	"time"

	"github.com/RangelReale/panyl/v2"
)

const PostgresLogFormat = "postgres_log"

// PostgresLog parses MongoDB log lines format
type PostgresLog struct {
}

var _ panyl.PluginParse = PostgresLog{}

// example: "2022-04-05 14:29:07.500 UTC [73] ERROR:  relation "users" does not exist at character 36"

var (
	PostgresLogRe           = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} [^\s]+)\sUTC\s\[(\d+)]\s(\w+):\s+(.*)`)
	postgresTimestampFormat = "2006-01-02 15:04:05.000"
)

func (m PostgresLog) ExtractParse(ctx context.Context, lines panyl.ItemLines, item *panyl.Item) (bool, error) {
	// Only single line is supported
	if len(lines) != 1 {
		return false, nil
	}

	matches := PostgresLogRe.FindStringSubmatch(lines.Line())
	if matches == nil {
		return false, nil
	}

	err := item.MergeLinesData(lines)
	if err != nil {
		return false, err
	}
	item.Line = ""

	timestamp := matches[1]
	pid := matches[2]
	level := matches[3]
	message := matches[4]

	item.Data["timestamp"] = timestamp
	item.Data["pid"] = pid
	item.Data["level"] = level
	item.Data["message"] = message

	item.Metadata[panyl.MetadataFormat] = PostgresLogFormat
	item.Metadata[panyl.MetadataMessage] = message

	if timestamp != "" {
		ts, err := time.Parse(postgresTimestampFormat, timestamp)
		if err == nil {
			item.Metadata[panyl.MetadataTimestamp] = ts
		}
	}

	// https://www.postgresql.org/docs/current/runtime-config-logging.html
	if level == "ERROR" || level == "FATAL" || level == "PANIC" {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelERROR
	} else if level == "WARNING" {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
	} else if level == "DEBUG" {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
	} else if level == "STATEMENT" || level == "DETAIL" {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelTRACE
	} else {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
	}

	return true, nil
}

func (m PostgresLog) IsPanylPlugin() {}
