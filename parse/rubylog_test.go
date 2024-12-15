package parse

import (
	"context"
	"testing"

	"github.com/RangelReale/panyl/v2"
	"github.com/stretchr/testify/assert"
)

func TestRubyLog(t *testing.T) {
	type test struct {
		source   string
		level    string
		category string
		message  string
	}

	tests := []test{
		{
			source:   `I, [1999-03-03T02:34:24.895701 #19074]  INFO -- Main: info.`,
			level:    panyl.MetadataLevelINFO,
			category: "",
			message:  "info.",
		},
	}

	for _, tc := range tests {
		ctx := context.Background()
		result := panyl.InitItem()

		p := &RubyLog{}
		ok, err := p.ExtractParse(ctx, panyl.ItemLines{&panyl.Item{Line: tc.source}}, result)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, result.Metadata[panyl.MetadataTimestamp])
		assert.Equal(t, tc.level, result.Metadata.StringValue(panyl.MetadataLevel))
		assert.Equal(t, tc.category, result.Metadata.StringValue(panyl.MetadataCategory))
		assert.Equal(t, tc.message, result.Metadata.StringValue(panyl.MetadataMessage))
	}
}
