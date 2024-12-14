package metadata

import (
	"github.com/RangelReale/panyl"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRubyForeman(t *testing.T) {
	result := panyl.InitProcess()
	result.Line = "16:41:59 api.1         | log text"

	plugin := &RubyForeman{}
	ok, err := plugin.ExtractMetadata(result)
	assert.NoError(t, err)
	assert.True(t, ok)

	assert.Equal(t, "api.1", result.Metadata.StringValue(panyl.MetadataApplication))
	assert.Equal(t, "log text", result.Line)
}
