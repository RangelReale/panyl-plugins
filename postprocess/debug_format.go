package postprocess

import (
	"context"
	"fmt"

	"github.com/RangelReale/panyl/v2"
)

var _ panyl.PluginPostProcess = (*DebugFormat)(nil)

type DebugFormat struct {
}

func (p DebugFormat) PostProcess(ctx context.Context, item *panyl.Item) (bool, error) {
	if message := item.Metadata.StringValue(panyl.MetadataMessage); message != "" {
		item.Metadata[panyl.MetadataMessage] = fmt.Sprintf("[[fmt:%s]] %s",
			item.Metadata.StringValue(panyl.MetadataFormat), message)
	}
	return false, nil
}

func (p DebugFormat) PostProcessOrder() int {
	return panyl.PostProcessOrderLast - 1
}

func (p DebugFormat) IsPanylPlugin() {}
