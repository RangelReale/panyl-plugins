package parse

import (
	"context"
	"regexp"
	"time"

	"github.com/RangelReale/panyl/v2"
)

const RedisLogFormat = "redis_log"

// RedisLog parses Redis log lines format
type RedisLog struct {
}

// example: "21:C 13 Apr 2022 17:59:51.096 * RDB: 0 MB of memory used by copy-on-write"

var _ panyl.PluginParse = RedisLog{}

var (
	RedisLogRe           = regexp.MustCompile(`^(\d+):(\w)\s+(\d{2}\s\w{3}\s\d{4} \d{2}:\d{2}:\d{2}.\d{3})\s([.*#-])\s*(.*)$`)
	redisTimestampFormat = "02 Jan 2006 15:04:05.000"
)

func (m RedisLog) ExtractParse(ctx context.Context, lines panyl.ItemLines, item *panyl.Item) (bool, error) {
	// https://github.com/redis/redis/issues/2545#issuecomment-97270522

	// Only single line is supported
	if len(lines) != 1 {
		return false, nil
	}

	matches := RedisLogRe.FindStringSubmatch(lines.Line())
	if matches == nil {
		return false, nil
	}

	err := item.MergeLinesData(lines)
	if err != nil {
		return false, err
	}
	item.Line = ""

	pid := matches[1]
	role := matches[2]
	timestamp := matches[3]
	level := matches[4]
	message := matches[5]

	item.Data["pid"] = pid
	item.Data["role"] = role
	item.Data["timestamp"] = timestamp
	item.Data["level"] = level
	item.Data["message"] = message

	item.Metadata[panyl.MetadataFormat] = RedisLogFormat
	item.Metadata[panyl.MetadataMessage] = message

	if timestamp != "" {
		ts, err := time.Parse(redisTimestampFormat, timestamp)
		if err == nil {
			item.Metadata[panyl.MetadataTimestamp] = ts
		}
	}

	// https://build47.com/redis-log-format-levels/
	switch level {
	case ".":
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
	case "-":
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelTRACE
	case "*":
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
	case "#":
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
	}

	switch role {
	case "X":
		item.Metadata[panyl.MetadataCategory] = "sentinel"
	case "C":
		item.Metadata[panyl.MetadataCategory] = "child"
	case "S":
		item.Metadata[panyl.MetadataCategory] = "slave"
	case "M":
		item.Metadata[panyl.MetadataCategory] = "master"
	}

	return true, nil
}

func (m RedisLog) IsPanylPlugin() {}
