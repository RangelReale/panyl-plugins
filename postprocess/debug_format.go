package postprocess

import (
	"fmt"

	"github.com/RangelReale/panyl"
)

var _ panyl.PluginPostProcess = (*DebugFormat)(nil)

type DebugFormat struct {
}

func (p DebugFormat) PostProcess(result *panyl.Process) (bool, error) {
	if message := result.Metadata.StringValue(panyl.Metadata_Message); message != "" {
		result.Metadata[panyl.Metadata_Message] = fmt.Sprintf("[[fmt:%s]] %s",
			result.Metadata.StringValue(panyl.Metadata_Format), message)
	}
	return false, nil
}

func (p DebugFormat) PostProcessOrder() int {
	return panyl.PostProcessOrder_Last - 1
}

func (p DebugFormat) IsPanylPlugin() {}
