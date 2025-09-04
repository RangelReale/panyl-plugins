package parse

import (
	"context"
	"testing"

	"github.com/RangelReale/panyl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErlangLog(t *testing.T) {
	type test struct {
		source   string
		level    string
		category string
		message  string
	}

	tests := []test{
		{
			source:  `2025-09-04T13:26:52.315705+00:00 [warning] not_all_kafka_partitions_connected`,
			level:   panyl.MetadataLevelWARNING,
			message: "not_all_kafka_partitions_connected",
		},
	}

	for _, tc := range tests {
		ctx := context.Background()
		item := panyl.InitItem()

		p := ErlangLog{}
		ok, err := p.ExtractParse(ctx, panyl.ItemLines{&panyl.Item{Line: tc.source}}, item)
		require.NoError(t, err)
		require.True(t, ok, "plugin didn't match")

		assert.NotZero(t, item.Metadata[panyl.MetadataTimestamp])
		assert.Equal(t, tc.level, item.Metadata.StringValue(panyl.MetadataLevel))
		assert.Equal(t, tc.message, item.Metadata.StringValue(panyl.MetadataMessage))
	}
}
