package metadata

import (
	"github.com/RangelReale/panyl"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDockerCompose(t *testing.T) {
	result := panyl.InitProcess()
	result.Line = "application    | my log here"

	plugin := &DockerCompose{}
	ok, err := plugin.ExtractMetadata(result)
	assert.NoError(t, err)
	assert.True(t, ok)

	assert.Equal(t, "application", result.Metadata.StringValue(panyl.Metadata_Application))
	assert.Equal(t, "my log here", result.Line)
}
