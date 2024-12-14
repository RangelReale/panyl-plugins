package parse

import (
	"github.com/RangelReale/panyl"
	"regexp"
	"time"
)

var _ panyl.PluginParse = (*PostgresLog)(nil)

const PostgresLog_Format = "postgres_log"

// PostgresLog parses MongoDB log lines format
type PostgresLog struct {
}

// example: "2022-04-05 14:29:07.500 UTC [73] ERROR:  relation "users" does not exist at character 36"

var (
	PostgresLogRe           = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} [^\s]+)\sUTC\s\[(\d+)]\s(\w+):\s+(.*)`)
	postgresTimestampFormat = "2006-01-02 15:04:05.000"
)

func (m *PostgresLog) ExtractParse(lines panyl.ProcessLines, result *panyl.Process) (bool, error) {
	// Only single line is supported
	if len(lines) != 1 {
		return false, nil
	}

	matches := PostgresLogRe.FindStringSubmatch(lines.Line())
	if matches == nil {
		return false, nil
	}

	err := result.MergeLinesData(lines)
	if err != nil {
		return false, err
	}
	result.Line = ""

	timestamp := matches[1]
	pid := matches[2]
	level := matches[3]
	message := matches[4]

	result.Data["timestamp"] = timestamp
	result.Data["pid"] = pid
	result.Data["level"] = level
	result.Data["message"] = message

	result.Metadata[panyl.MetadataFormat] = PostgresLog_Format
	result.Metadata[panyl.MetadataMessage] = message

	if timestamp != "" {
		ts, err := time.Parse(postgresTimestampFormat, timestamp)
		if err == nil {
			result.Metadata[panyl.MetadataTimestamp] = ts
		}
	}

	// https://www.postgresql.org/docs/current/runtime-config-logging.html
	if level == "ERROR" {
		result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelERROR
	} else if level == "FATAL" || level == "PANIC" {
		result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevel_FATAL
	} else if level == "WARNING" {
		result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
	} else if level == "DEBUG" {
		result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
	} else if level == "STATEMENT" || level == "DETAIL" {
		result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelTRACE
	} else {
		result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
	}

	return true, nil
}

func (m PostgresLog) IsPanylPlugin() {}
