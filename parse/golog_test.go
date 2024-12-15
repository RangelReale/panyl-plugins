package parse

import (
	"context"
	"testing"

	"github.com/RangelReale/panyl/v2"
	"github.com/stretchr/testify/assert"
)

func TestGoLog(t *testing.T) {
	type test struct {
		source   string
		level    string
		category string
		message  string
	}

	tests := []test{
		{
			source:   `2022-03-10T19:53:21.434Z	INFO	datadog-go/tracer.go:35	Datadog Tracer v1.28.0 ERROR: lost 2 traces`,
			level:    panyl.MetadataLevelINFO,
			category: "datadog-go/tracer",
			message:  "Datadog Tracer v1.28.0 ERROR: lost 2 traces",
		},
	}

	for _, tc := range tests {
		ctx := context.Background()

		result := panyl.InitItem()

		p := &GoLog{SourceAsCategory: true}
		ok, err := p.ExtractParse(ctx, panyl.ItemLines{&panyl.Item{Line: tc.source}}, result)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, result.Metadata[panyl.MetadataTimestamp])
		assert.Equal(t, tc.level, result.Metadata.StringValue(panyl.MetadataLevel))
		assert.Equal(t, tc.category, result.Metadata.StringValue(panyl.MetadataCategory))
		assert.Equal(t, tc.message, result.Metadata.StringValue(panyl.MetadataMessage))
	}
}
