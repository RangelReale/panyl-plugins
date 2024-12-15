package parse

import (
	"context"
	"regexp"
	"time"

	"github.com/RangelReale/panyl/v2"
)

const MongoLogFormat = "mongo_log"

// MongoLog parses MongoDB log lines format
type MongoLog struct {
}

// example: "2022-03-10T20:17:25.039+0000 I  NETWORK  [conn15] end connection 172.19.0.1:49812 (7 connections now open)"

var _ panyl.PluginParse = MongoLog{}

var (
	mongoLogRe           = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T[^\s]+)\s+(\w)\s+(\w+)\s+\[(\w+)]\s+(.*)$`)
	mongoTimestampFormat = "2006-01-02T15:04:05.999999999-0700"
)

func (m MongoLog) ExtractParse(ctx context.Context, lines panyl.ItemLines, item *panyl.Item) (bool, error) {
	// Only single line is supported
	if len(lines) != 1 {
		return false, nil
	}

	matches := mongoLogRe.FindStringSubmatch(lines.Line())
	if matches == nil {
		return false, nil
	}

	err := item.MergeLinesData(lines)
	if err != nil {
		return false, err
	}
	item.Line = ""

	timestamp := matches[1]
	severity := matches[2]
	component := matches[3]
	message := matches[5]

	item.Data["timestamp"] = timestamp
	item.Data["severity"] = severity
	item.Data["component"] = component
	item.Data["context"] = matches[4]
	item.Data["message"] = message

	item.Metadata[panyl.MetadataFormat] = MongoLogFormat
	item.Metadata[panyl.MetadataMessage] = message
	item.Metadata[panyl.MetadataCategory] = component

	if timestamp != "" {
		ts, err := time.Parse(mongoTimestampFormat, timestamp)
		if err == nil {
			item.Metadata[panyl.MetadataTimestamp] = ts
		}
	}

	// https://docs.mongodb.com/manual/reference/log-messages/
	if severity == "D" {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
	} else if severity == "I" {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
	} else if severity == "W" {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
	} else if severity == "E" {
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelERROR
	}

	return true, nil
}

func (m MongoLog) IsPanylPlugin() {}
