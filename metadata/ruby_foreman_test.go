package metadata

import (
	"context"
	"github.com/RangelReale/panyl/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRubyForeman(t *testing.T) {
	ctx := context.Background()

	result := panyl.InitProcess()
	result.Line = "16:41:59 api.1         | log text"

	plugin := &RubyForeman{}
	ok, err := plugin.ExtractMetadata(ctx, result)
	assert.NoError(t, err)
	assert.True(t, ok)

	assert.Equal(t, "api.1", result.Metadata.StringValue(panyl.MetadataApplication))
	assert.Equal(t, "log text", result.Line)
}
