package parseformat

import (
	"github.com/RangelReale/panyl"
	"time"
)

var _ panyl.PluginParseFormat = (*ElasticSearchJSON)(nil)

const ElasticSearchJSONFormat = "elasticsearch_json"

type ElasticSearchJSON struct {
}

// example: {"cluster.name":"docker-cluster","component":"o.e.t.TransportService","level":"INFO","message":"publish_address {172.18.0.4:9300}, bound_addresses {0.0.0.0:9300}","node.name":"3404ffa7b26c","timestamp":"2022-04-13T17:24:56,134Z","type":"server"}

var (
	elasticSearchTimestampFormat = "2006-01-02T15:04:05,000Z07:00"
)

func (C ElasticSearchJSON) ParseFormat(result *panyl.Process) (bool, error) {
	// only if json
	if result.Metadata.StringValue(panyl.Metadata_Structure) == panyl.MetadataStructure_JSON {
		if result.Data.HasValue("timestamp") && result.Data.HasValue("cluster.name") &&
			result.Data.HasValue("node.name") && result.Data.HasValue("type") {
			timestamp := result.Data.StringValue("timestamp")
			level := result.Data.StringValue("level")
			message := result.Data.StringValue("message")
			//component := result.Data.StringValue("component")
			typ := result.Data.StringValue("type")

			result.Metadata[panyl.Metadata_Format] = ElasticSearchJSONFormat
			result.Metadata[panyl.Metadata_Message] = message
			result.Metadata[panyl.Metadata_Category] = typ

			if timestamp != "" {
				ts, err := time.Parse(elasticSearchTimestampFormat, timestamp)
				if err == nil {
					result.Metadata[panyl.Metadata_Timestamp] = ts
				}
			}

			// https://www.elastic.co/guide/en/elasticsearch/reference/current/logging.html
			switch level {
			case "OFF", "FATAL":
				result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_FATAL
			case "ERROR":
				result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_ERROR
			case "WARN":
				result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_WARNING
			case "INFO":
				result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_INFO
			case "DEBUG":
				result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_DEBUG
			case "TRACE":
				result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_TRACE
			}
			return true, nil
		}
	}
	return false, nil
}

func (C ElasticSearchJSON) IsPanylPlugin() {}
