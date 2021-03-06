package parse

import (
	"github.com/RangelReale/panyl"
	"regexp"
	"time"
)

var _ panyl.PluginParse = (*RedisLog)(nil)

const RedisLog_Format = "redis_log"

// RedisLog parses Redis log lines format
type RedisLog struct {
}

// example: "21:C 13 Apr 2022 17:59:51.096 * RDB: 0 MB of memory used by copy-on-write"

var (
	RedisLogRe           = regexp.MustCompile(`^(\d+):(\w)\s+(\d{2}\s\w{3}\s\d{4} \d{2}:\d{2}:\d{2}.\d{3})\s([.*#-])\s*(.*)$`)
	redisTimestampFormat = "02 Jan 2006 15:04:05.000"
)

// https://github.com/redis/redis/issues/2545#issuecomment-97270522
func (m *RedisLog) ExtractParse(lines panyl.ProcessLines, result *panyl.Process) (bool, error) {
	// Only single line is supported
	if len(lines) != 1 {
		return false, nil
	}

	matches := RedisLogRe.FindStringSubmatch(lines.Line())
	if matches == nil {
		return false, nil
	}

	err := result.MergeLinesData(lines)
	if err != nil {
		return false, err
	}
	result.Line = ""

	pid := matches[1]
	role := matches[2]
	timestamp := matches[3]
	level := matches[4]
	message := matches[5]

	result.Data["pid"] = pid
	result.Data["role"] = role
	result.Data["timestamp"] = timestamp
	result.Data["level"] = level
	result.Data["message"] = message

	result.Metadata[panyl.Metadata_Format] = RedisLog_Format
	result.Metadata[panyl.Metadata_Message] = message

	if timestamp != "" {
		ts, err := time.Parse(redisTimestampFormat, timestamp)
		if err == nil {
			result.Metadata[panyl.Metadata_Timestamp] = ts
		}
	}

	// https://build47.com/redis-log-format-levels/
	switch level {
	case ".":
		result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_DEBUG
	case "-":
		result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_TRACE
	case "*":
		result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_INFO
	case "#":
		result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_WARNING
	}

	switch role {
	case "X":
		result.Metadata[panyl.Metadata_Category] = "sentinel"
	case "C":
		result.Metadata[panyl.Metadata_Category] = "child"
	case "S":
		result.Metadata[panyl.Metadata_Category] = "slave"
	case "M":
		result.Metadata[panyl.Metadata_Category] = "master"
	}

	return true, nil
}

func (m RedisLog) IsPanylPlugin() {}
