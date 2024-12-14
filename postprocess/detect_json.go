package postprocess

import (
	"context"
	"strings"
	"time"

	"github.com/RangelReale/panyl/v2"
)

var _ panyl.PluginPostProcess = (*DetectJSON)(nil)

type DetectJSON struct{}

func (p DetectJSON) PostProcess(ctx context.Context, result *panyl.Process) (bool, error) {
	if result.Metadata.HasValue(panyl.MetadataFormat) {
		// already has a known format
		return false, nil
	}

	// only if json
	if result.Metadata.StringValue(panyl.MetadataStructure) == panyl.MetadataStructureJSON {
		// timestamp
		if !result.Metadata.HasValue(panyl.MetadataTimestamp) {
			var detectTimestamp string
			if result.Data.HasValue("timestamp") {
				detectTimestamp = result.Data.StringValue("timestamp")
			} else if result.Data.HasValue("time") {
				detectTimestamp = result.Data.StringValue("time")
			}
			if detectTimestamp != "" {
				if ts, err := time.Parse(time.RFC3339, detectTimestamp); err != nil {
					result.Metadata[panyl.MetadataTimestamp] = ts
				} else if ts, err := time.Parse(time.RFC3339Nano, detectTimestamp); err != nil {
					result.Metadata[panyl.MetadataTimestamp] = ts
				}
			}
		}

		// level
		if !result.Metadata.HasValue(panyl.MetadataLevel) {
			var detectLevel string
			if result.Data.HasValue("level") {
				detectLevel = result.Data.StringValue("level")
			}
			if detectLevel != "" {
				switch strings.ToLower(detectLevel) {
				case "error", "fatal":
					result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelERROR
				case "warn", "warning":
					result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
				case "info", "information":
					result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
				case "debug":
					result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
				case "trace":
					result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelTRACE
				}
			}
		}

		// level
		if !result.Metadata.HasValue(panyl.MetadataMessage) {
			var detectMessage string
			if result.Data.HasValue("message") {
				detectMessage = result.Data.StringValue("message")
			}
			if detectMessage != "" {
				result.Metadata[panyl.MetadataMessage] = detectMessage
			}
		}
	}
	return false, nil
}

func (p DetectJSON) PostProcessOrder() int {
	return panyl.PostProcessOrderFirst + 1
}

func (p DetectJSON) IsPanylPlugin() {}
