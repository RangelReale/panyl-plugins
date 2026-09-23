package parse

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/RangelReale/panyl/v2"
)

const PostgresLogFormat = "postgres_log"

// PostgresLog parses Postgres log lines format
type PostgresLog struct {
	// Location is the time zone used to parse timestamps whose zone is an abbreviation other than UTC/GMT
	// (e.g. "CET"), as abbreviations are ambiguous. If nil, the abbreviation is resolved by [time.Parse].
	Location *time.Location
}

var _ panyl.PluginParse = PostgresLog{}

// example: "2022-04-05 14:29:07.500 UTC [73] ERROR:  relation "users" does not exist at character 36"
// example: "2022-04-05 11:29:07 -03 [73] LOG:  database system is ready to accept connections"

var (
	PostgresLogRe           = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(?:\.\d+)?)\s([A-Za-z]+|[+-]\d{2}(?::?\d{2})?)\s\[(\d+)]\s(\w+):\s+(.*)`)
	postgresTimestampFormat = "2006-01-02 15:04:05"
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
	timezone := matches[2]
	pid := matches[3]
	level := matches[4]
	message := matches[5]

	item.Data["timestamp"] = timestamp
	item.Data["timezone"] = timezone
	item.Data["pid"] = pid
	item.Data["level"] = level
	item.Data["message"] = message

	item.Metadata[panyl.MetadataFormat] = PostgresLogFormat
	item.Metadata[panyl.MetadataMessage] = message

	if timestamp != "" {
		ts, err := m.parseTimestamp(timestamp, timezone)
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

func (m PostgresLog) parseTimestamp(timestamp, timezone string) (time.Time, error) {
	switch {
	case timezone == "UTC" || timezone == "GMT":
		return time.ParseInLocation(postgresTimestampFormat, timestamp, time.UTC)
	case timezone[0] == '+' || timezone[0] == '-':
		// numeric offset: "-03", "+0530" or "+05:30"
		for _, zoneLayout := range []string{"-07", "-0700", "-07:00"} {
			if ts, err := time.Parse(postgresTimestampFormat+" "+zoneLayout, timestamp+" "+timezone); err == nil {
				return ts, nil
			}
		}
		return time.Time{}, fmt.Errorf("invalid time zone offset %q", timezone)
	case m.Location != nil:
		return time.ParseInLocation(postgresTimestampFormat, timestamp, m.Location)
	default:
		return time.Parse(postgresTimestampFormat+" MST", timestamp+" "+timezone)
	}
}
