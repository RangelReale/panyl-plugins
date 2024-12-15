package postprocess

import (
	"context"
	"fmt"

	"github.com/RangelReale/panyl/v2"
)

type DebugFormat struct {
}

var _ panyl.PluginPostProcess = DebugFormat{}

func (m DebugFormat) PostProcess(ctx context.Context, item *panyl.Item) (bool, error) {
	if message := item.Metadata.StringValue(panyl.MetadataMessage); message != "" {
		item.Metadata[panyl.MetadataMessage] = fmt.Sprintf("[[fmt:%s]] %s",
			item.Metadata.StringValue(panyl.MetadataFormat), message)
	}
	return false, nil
}

func (m DebugFormat) PostProcessOrder() int {
	return panyl.PostProcessOrderLast - 1
}

func (m DebugFormat) IsPanylPlugin() {}
