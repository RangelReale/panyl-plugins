package parse

import (
	"context"
	"regexp"
	"time"

	"github.com/RangelReale/panyl/v2"
)

const JavaLogFormat = "java_log"

// JavaLog parse Golang log lines format
type JavaLog struct{}

var _ panyl.PluginParse = JavaLog{}

// example: "2025-02-05 19:44:48,878 INFO Processing ruok command from /127.0.0.1:48544 (org.apache.zookeeper.server.NettyServerCnxn) [nioEventLoopGroup-4-3]"

var javaLogRe = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} [^\s]+)\s+(\w+)\s+(.*)$`)

const javaTimeFormat = "2006-01-02 15:04:05,999"

func (m JavaLog) ExtractParse(ctx context.Context, lines panyl.ItemLines, item *panyl.Item) (bool, error) {
	// Only single line is supported
	if len(lines) != 1 {
		return false, nil
	}

	matches := javaLogRe.FindStringSubmatch(lines.Line())
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

	item.Metadata[panyl.MetadataFormat] = JavaLogFormat
	item.Metadata[panyl.MetadataMessage] = message

	if timestamp != "" {
		ts, err := time.Parse(javaTimeFormat, timestamp)
		if err == nil {
			item.Metadata[panyl.MetadataTimestamp] = ts
		}
	}
	if level == "DEBUG" {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
	} else if level == "INFO" {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
	} else if level == "WARN" {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
	} else if level == "ERROR" {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelERROR
	}

	return true, nil
}

func (m JavaLog) IsPanylPlugin() {}
