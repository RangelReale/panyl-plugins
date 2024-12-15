package postprocess

import (
	"context"
	"strings"
	"time"

	"github.com/RangelReale/panyl/v2"
)

var _ panyl.PluginPostProcess = (*DetectJSON)(nil)

type DetectJSON struct{}

func (p DetectJSON) PostProcess(ctx context.Context, item *panyl.Item) (bool, error) {
	if item.Metadata.HasValue(panyl.MetadataFormat) {
		// already has a known format
		return false, nil
	}

	// only if json
	if item.Metadata.StringValue(panyl.MetadataStructure) == panyl.MetadataStructureJSON {
		// timestamp
		if !item.Metadata.HasValue(panyl.MetadataTimestamp) {
			var detectTimestamp string
			if item.Data.HasValue("timestamp") {
				detectTimestamp = item.Data.StringValue("timestamp")
			} else if item.Data.HasValue("time") {
				detectTimestamp = item.Data.StringValue("time")
			}
			if detectTimestamp != "" {
				if ts, err := time.Parse(time.RFC3339, detectTimestamp); err != nil {
					item.Metadata[panyl.MetadataTimestamp] = ts
				} else if ts, err := time.Parse(time.RFC3339Nano, detectTimestamp); err != nil {
					item.Metadata[panyl.MetadataTimestamp] = ts
				}
			}
		}

		// level
		if !item.Metadata.HasValue(panyl.MetadataLevel) {
			var detectLevel string
			if item.Data.HasValue("level") {
				detectLevel = item.Data.StringValue("level")
			}
			if detectLevel != "" {
				switch strings.ToLower(detectLevel) {
				case "error", "fatal":
					item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelERROR
				case "warn", "warning":
					item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
				case "info", "information":
					item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
				case "debug":
					item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
				case "trace":
					item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelTRACE
				}
			}
		}

		// level
		if !item.Metadata.HasValue(panyl.MetadataMessage) {
			var detectMessage string
			if item.Data.HasValue("message") {
				detectMessage = item.Data.StringValue("message")
			}
			if detectMessage != "" {
				item.Metadata[panyl.MetadataMessage] = detectMessage
			}
		}
	}
	return false, nil
}

func (p DetectJSON) PostProcessOrder() int {
	return panyl.PostProcessOrderFirst + 1
}

func (p DetectJSON) IsPanylPlugin() {}
