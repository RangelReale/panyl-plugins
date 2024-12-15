package postprocess

import (
	"context"
	"fmt"

	"github.com/RangelReale/panyl/v2"
)

var _ panyl.PluginPostProcess = (*DebugFormat)(nil)

type DebugFormat struct {
}

func (p DebugFormat) PostProcess(ctx context.Context, result *panyl.Item) (bool, error) {
	if message := result.Metadata.StringValue(panyl.MetadataMessage); message != "" {
		result.Metadata[panyl.MetadataMessage] = fmt.Sprintf("[[fmt:%s]] %s",
			result.Metadata.StringValue(panyl.MetadataFormat), message)
	}
	return false, nil
}

func (p DebugFormat) PostProcessOrder() int {
	return panyl.PostProcessOrderLast - 1
}

func (p DebugFormat) IsPanylPlugin() {}
