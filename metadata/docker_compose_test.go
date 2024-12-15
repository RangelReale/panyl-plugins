package metadata

import (
	"context"
	"testing"

	"github.com/RangelReale/panyl/v2"
	"github.com/stretchr/testify/assert"
)

func TestDockerCompose(t *testing.T) {
	ctx := context.Background()

	item := panyl.InitItem()
	item.Line = "application    | my log here"

	plugin := &DockerCompose{}
	ok, err := plugin.ExtractMetadata(ctx, item)
	assert.NoError(t, err)
	assert.True(t, ok)

	assert.Equal(t, "application", item.Metadata.StringValue(panyl.MetadataApplication))
	assert.Equal(t, "my log here", item.Line)
}
