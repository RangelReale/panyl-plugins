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

	plugin := DockerCompose{}
	ok, err := plugin.ExtractMetadata(ctx, item)
	assert.NoError(t, err)
	assert.True(t, ok)

	assert.Equal(t, "application", item.Metadata.StringValue(panyl.MetadataApplication))
	assert.Equal(t, "my log here", item.Line)
}

func TestDockerComposeSeparator(t *testing.T) {
	tests := []struct {
		source      string
		application string
		line        string
	}{
		{source: "app |message", application: "app", line: "message"},
		{source: "app-1  |  indented", application: "app-1", line: " indented"},
		{source: "app |", application: "app", line: ""},
	}

	for _, tc := range tests {
		ctx := context.Background()
		item := panyl.InitItem()
		item.Line = tc.source

		ok, err := DockerCompose{}.ExtractMetadata(ctx, item)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.Equal(t, tc.application, item.Metadata.StringValue(panyl.MetadataApplication))
		assert.Equal(t, tc.line, item.Line)
	}
}
