package metadata

import (
	"context"
	"testing"

	"github.com/RangelReale/panyl/v2"
	"github.com/stretchr/testify/assert"
)

func TestRubyForeman(t *testing.T) {
	ctx := context.Background()

	item := panyl.InitItem()
	item.Line = "16:41:59 api.1         | log text"

	plugin := &RubyForeman{}
	ok, err := plugin.ExtractMetadata(ctx, item)
	assert.NoError(t, err)
	assert.True(t, ok)

	assert.Equal(t, "api.1", item.Metadata.StringValue(panyl.MetadataApplication))
	assert.Equal(t, "log text", item.Line)
}
