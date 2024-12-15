package parse

import (
	"context"
	"testing"

	"github.com/RangelReale/panyl/v2"
	"github.com/stretchr/testify/assert"
)

func TestPostgresLog(t *testing.T) {
	type test struct {
		source   string
		level    string
		category string
		message  string
	}

	tests := []test{
		{
			source:   `2022-04-05 14:29:07.500 UTC [73] ERROR:  relation "users" does not exist at character 36`,
			level:    panyl.MetadataLevelERROR,
			category: "",
			message:  `relation "users" does not exist at character 36`,
		},
	}

	for _, tc := range tests {
		ctx := context.Background()
		result := panyl.InitItem()

		p := &PostgresLog{}
		ok, err := p.ExtractParse(ctx, panyl.ItemLines{&panyl.Item{Line: tc.source}}, result)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, result.Metadata[panyl.MetadataTimestamp])
		assert.Equal(t, tc.level, result.Metadata.StringValue(panyl.MetadataLevel))
		assert.Equal(t, tc.category, result.Metadata.StringValue(panyl.MetadataCategory))
		assert.Equal(t, tc.message, result.Metadata.StringValue(panyl.MetadataMessage))
	}
}
