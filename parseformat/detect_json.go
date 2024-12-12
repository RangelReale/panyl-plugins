package parseformat

import (
	"strings"
	"time"

	"github.com/RangelReale/panyl"
)

var _ panyl.PluginParseFormat = (*DetectJSON)(nil)

type DetectJSON struct{}

// example: {"cluster.name":"docker-cluster","component":"o.e.t.TransportService","level":"INFO","message":"publish_address {172.18.0.4:9300}, bound_addresses {0.0.0.0:9300}","node.name":"3404ffa7b26c","timestamp":"2022-04-13T17:24:56,134Z","type":"server"}

func (C DetectJSON) ParseFormat(result *panyl.Process) (bool, error) {
	// only if json
	if result.Metadata.StringValue(panyl.Metadata_Structure) == panyl.MetadataStructure_JSON {
		// timestamp
		if !result.Metadata.HasValue(panyl.Metadata_Timestamp) {
			var detectTimestamp string
			if result.Data.HasValue("timestamp") {
				detectTimestamp = result.Data.StringValue("timestamp")
			} else if result.Data.HasValue("time") {
				detectTimestamp = result.Data.StringValue("time")
			}
			if detectTimestamp != "" {
				if ts, err := time.Parse(time.RFC3339, detectTimestamp); err != nil {
					result.Metadata[panyl.Metadata_Timestamp] = ts
				} else if ts, err := time.Parse(time.RFC3339Nano, detectTimestamp); err != nil {
					result.Metadata[panyl.Metadata_Timestamp] = ts
				}
			}
		}

		// level
		if !result.Metadata.HasValue(panyl.Metadata_Level) {
			var detectLevel string
			if result.Data.HasValue("level") {
				detectLevel = result.Data.StringValue("level")
			}
			if detectLevel != "" {
				switch strings.ToLower(detectLevel) {
				case "fatal":
					result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_FATAL
				case "error":
					result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_ERROR
				case "warn", "warning":
					result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_WARNING
				case "info", "information":
					result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_INFO
				case "debug":
					result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_DEBUG
				case "trace":
					result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_TRACE
				}
			}
		}

		// level
		if !result.Metadata.HasValue(panyl.Metadata_Message) {
			var detectMessage string
			if result.Data.HasValue("message") {
				detectMessage = result.Data.StringValue("message")
			}
			if detectMessage != "" {
				result.Metadata[panyl.Metadata_Message] = detectMessage
			}
		}
	}
	return false, nil
}

func (C DetectJSON) IsPanylPlugin() {}
