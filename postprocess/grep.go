package postprocess

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/RangelReale/panyl/v2"
)

type Grep struct {
	Values []string
}

var _ panyl.PluginPostProcess = (*Grep)(nil)

func (g Grep) IsPanylPlugin() {}

func (g Grep) PostProcessOrder() int {
	return panyl.PostProcessOrderFirst + 1
}

func (g Grep) PostProcess(ctx context.Context, item *panyl.Item) (bool, error) {
	for _, value := range g.Values {
		var message string
		if msg := item.Metadata.StringValue(panyl.MetadataMessage); msg != "" {
			message = msg
		} else if len(item.Data) > 0 {
			dt, err := json.Marshal(item.Data)
			if err != nil {
				message = fmt.Sprintf("Error marshaling data to json: %s", err.Error())
			} else {
				message = string(dt)
			}
		} else if item.Line != "" {
			message = item.Line
		}

		if message != "" {
			if strings.Contains(strings.ToLower(message), strings.ToLower(strings.TrimSpace(value))) {
				newItem, err := item.Clone()
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "unexpected error cloning item: %v\n", err)
					continue
				}
				category := newItem.Metadata.StringValue(panyl.MetadataCategory)
				if category != "" {
					category = fmt.Sprintf("[grep](%s)[%s]", category, value)
				} else {
					category = fmt.Sprintf("[grep][%s]", value)
				}
				item.Metadata.ListValueAdd(panyl.MetadataExtraCategories, category)
			}
		}
	}
	return true, nil
}
