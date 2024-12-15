package parseformat

import (
	"context"
	"time"

	"github.com/RangelReale/panyl/v2"
)

var _ panyl.PluginParseFormat = (*ElasticSearchJSON)(nil)

const ElasticSearchJSONFormat = "elasticsearch_json"

type ElasticSearchJSON struct {
}

// example: {"cluster.name":"docker-cluster","component":"o.e.t.TransportService","level":"INFO","message":"publish_address {172.18.0.4:9300}, bound_addresses {0.0.0.0:9300}","node.name":"3404ffa7b26c","timestamp":"2022-04-13T17:24:56,134Z","type":"server"}

var (
	elasticSearchTimestampFormat = "2006-01-02T15:04:05,000Z07:00"
)

func (C ElasticSearchJSON) ParseFormat(ctx context.Context, result *panyl.Item) (bool, error) {
	// only if json
	if result.Metadata.StringValue(panyl.MetadataStructure) == panyl.MetadataStructureJSON {
		if result.Data.HasValue("timestamp") && result.Data.HasValue("cluster.name") &&
			result.Data.HasValue("node.name") && result.Data.HasValue("type") {
			timestamp := result.Data.StringValue("timestamp")
			level := result.Data.StringValue("level")
			message := result.Data.StringValue("message")
			// component := result.Data.StringValue("component")
			typ := result.Data.StringValue("type")

			result.Metadata[panyl.MetadataFormat] = ElasticSearchJSONFormat
			result.Metadata[panyl.MetadataMessage] = message
			result.Metadata[panyl.MetadataCategory] = typ

			if timestamp != "" {
				ts, err := time.Parse(elasticSearchTimestampFormat, timestamp)
				if err == nil {
					result.Metadata[panyl.MetadataTimestamp] = ts
				}
			}

			// https://www.elastic.co/guide/en/elasticsearch/reference/current/logging.html
			switch level {
			case "ERROR", "OFF", "FATAL":
				result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelERROR
			case "WARN":
				result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
			case "INFO":
				result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
			case "DEBUG":
				result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
			case "TRACE":
				result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelTRACE
			}
			return true, nil
		}
	}
	return false, nil
}

func (C ElasticSearchJSON) IsPanylPlugin() {}
