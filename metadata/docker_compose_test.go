package metadata

import (
	"context"
	"testing"

	"github.com/RangelReale/panyl"
	"github.com/stretchr/testify/assert"
)

func TestDockerCompose(t *testing.T) {
	ctx := context.Background()

	result := panyl.InitProcess()
	result.Line = "application    | my log here"

	plugin := &DockerCompose{}
	ok, err := plugin.ExtractMetadata(ctx, result)
	assert.NoError(t, err)
	assert.True(t, ok)

	assert.Equal(t, "application", result.Metadata.StringValue(panyl.MetadataApplication))
	assert.Equal(t, "my log here", result.Line)
}
