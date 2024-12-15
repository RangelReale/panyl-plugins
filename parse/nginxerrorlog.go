package parse

import (
	"context"
	"regexp"
	"time"

	"github.com/RangelReale/panyl/v2"
)

var _ panyl.PluginParse = (*NGINXErrorLog)(nil)

const NGINXErrorLog_Format = "nginx_error_log"

// NGINXErrorLog parses NGINX log lines format
type NGINXErrorLog struct {
}

// 2022/03/10 20:20:48 [error] 62#62: *70 invalid URL prefix in "", client: 127.0.0.1, server: , request: "GET / HTTP/1.1", host: "localhost:8080", referrer: "http://localhost:8080"

var (
	// https://stackoverflow.com/a/26125951/784175
	nginxErrorLogRe           = regexp.MustCompile(`^(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2})\s+\[(\w+)]\s+(\d+)#(\d+):\s+((\*\w+)\s+)?(.*)$`)
	nginxErrorTimestampFormat = "2006/01/02 15:04:05"
)

func (m *NGINXErrorLog) ExtractParse(ctx context.Context, lines panyl.ItemLines, item *panyl.Item) (bool, error) {
	// Only single line is supported
	if len(lines) != 1 {
		return false, nil
	}

	matches := nginxErrorLogRe.FindStringSubmatch(lines.Line())
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
	message := matches[7]

	item.Data["timestamp"] = timestamp
	item.Data["level"] = level
	item.Data["pid"] = matches[3]
	item.Data["tid"] = matches[4]
	item.Data["cid"] = matches[6]
	item.Data["message"] = message

	item.Metadata[panyl.MetadataFormat] = NGINXErrorLog_Format
	item.Metadata[panyl.MetadataMessage] = message

	if timestamp != "" {
		ts, err := time.Parse(nginxErrorTimestampFormat, timestamp)
		if err == nil {
			item.Metadata[panyl.MetadataTimestamp] = ts
		}
	}

	// https://github.com/phusion/nginx/blob/master/src/core/ngx_log.c#L56
	switch level {
	case "debug":
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
	case "info", "notice":
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
	case "warn":
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
	case "error", "alert", "crit", "emerg":
		item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelERROR
	}

	return true, nil
}

func (m NGINXErrorLog) IsPanylPlugin() {}
