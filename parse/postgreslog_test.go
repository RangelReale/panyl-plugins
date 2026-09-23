package parse

import (
	"context"
	"testing"
	"time"

	"github.com/RangelReale/panyl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		{
			source:  `2022-04-05 11:29:07 -03 [73] LOG:  database system is ready to accept connections`,
			level:   panyl.MetadataLevelINFO,
			message: "database system is ready to accept connections",
		},
		{
			source:  `2022-04-05 19:59:07.5 +0530 [73] WARNING:  out of shared memory`,
			level:   panyl.MetadataLevelWARNING,
			message: "out of shared memory",
		},
	}

	for _, tc := range tests {
		ctx := context.Background()
		item := panyl.InitItem()

		p := PostgresLog{}
		ok, err := p.ExtractParse(ctx, panyl.ItemLines{&panyl.Item{Line: tc.source}}, item)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, item.Metadata[panyl.MetadataTimestamp])
		assert.Equal(t, tc.level, item.Metadata.StringValue(panyl.MetadataLevel))
		assert.Equal(t, tc.category, item.Metadata.StringValue(panyl.MetadataCategory))
		assert.Equal(t, tc.message, item.Metadata.StringValue(panyl.MetadataMessage))
	}
}

func TestPostgresLogTimezone(t *testing.T) {
	expected := time.Date(2022, 4, 5, 14, 29, 7, 0, time.UTC)

	tests := []struct {
		source   string
		location *time.Location
	}{
		{source: `2022-04-05 14:29:07 UTC [73] LOG:  x`},
		{source: `2022-04-05 14:29:07 GMT [73] LOG:  x`},
		{source: `2022-04-05 11:29:07 -03 [73] LOG:  x`},
		{source: `2022-04-05 19:59:07 +0530 [73] LOG:  x`},
		{source: `2022-04-05 19:59:07 +05:30 [73] LOG:  x`},
		{source: `2022-04-05 16:29:07 CEST [73] LOG:  x`, location: time.FixedZone("CEST", 2*60*60)},
	}

	for _, tc := range tests {
		ctx := context.Background()
		item := panyl.InitItem()

		ok, err := PostgresLog{Location: tc.location}.ExtractParse(ctx, panyl.ItemLines{&panyl.Item{Line: tc.source}}, item)
		require.NoError(t, err)
		require.True(t, ok, tc.source)

		ts, _ := item.Metadata[panyl.MetadataTimestamp].(time.Time)
		assert.True(t, expected.Equal(ts), "%s: expected %s, got %s", tc.source, expected, ts)
	}
}
