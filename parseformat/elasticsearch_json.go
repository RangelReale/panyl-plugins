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

func (C ElasticSearchJSON) ParseFormat(ctx context.Context, item *panyl.Item) (bool, error) {
	// only if json
	if item.Metadata.StringValue(panyl.MetadataStructure) == panyl.MetadataStructureJSON {
		if item.Data.HasValue("timestamp") && item.Data.HasValue("cluster.name") &&
			item.Data.HasValue("node.name") && item.Data.HasValue("type") {
			timestamp := item.Data.StringValue("timestamp")
			level := item.Data.StringValue("level")
			message := item.Data.StringValue("message")
			// component := item.Data.StringValue("component")
			typ := item.Data.StringValue("type")

			item.Metadata[panyl.MetadataFormat] = ElasticSearchJSONFormat
			item.Metadata[panyl.MetadataMessage] = message
			item.Metadata[panyl.MetadataCategory] = typ

			if timestamp != "" {
				ts, err := time.Parse(elasticSearchTimestampFormat, timestamp)
				if err == nil {
					item.Metadata[panyl.MetadataTimestamp] = ts
				}
			}

			// https://www.elastic.co/guide/en/elasticsearch/reference/current/logging.html
			switch level {
			case "ERROR", "OFF", "FATAL":
				item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelERROR
			case "WARN":
				item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
			case "INFO":
				item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
			case "DEBUG":
				item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
			case "TRACE":
				item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelTRACE
			}
			return true, nil
		}
	}
	return false, nil
}

func (C ElasticSearchJSON) IsPanylPlugin() {}
