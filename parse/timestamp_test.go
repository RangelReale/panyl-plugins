package parse

import (
	"context"
	"testing"
	"time"

	"github.com/RangelReale/panyl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimestampLocation(t *testing.T) {
	loc := time.FixedZone("TEST", -3*60*60)

	tests := []struct {
		name     string
		plugin   func(loc *time.Location) panyl.PluginParse
		source   string
		expected time.Time
	}{
		{
			name:     "ruby",
			plugin:   func(loc *time.Location) panyl.PluginParse { return RubyLog{Location: loc} },
			source:   "I, [1999-03-03T02:34:24.895701 #19074]  INFO -- Main: info.",
			expected: time.Date(1999, 3, 3, 2, 34, 24, 895701000, loc),
		},
		{
			name:     "java",
			plugin:   func(loc *time.Location) panyl.PluginParse { return JavaLog{Location: loc} },
			source:   "2025-02-05 19:44:48,878 INFO Processing",
			expected: time.Date(2025, 2, 5, 19, 44, 48, 878000000, loc),
		},
		{
			name:     "redis",
			plugin:   func(loc *time.Location) panyl.PluginParse { return RedisLog{Location: loc} },
			source:   "21:C 13 Apr 2022 17:59:51.096 * RDB: 0 MB of memory used by copy-on-write",
			expected: time.Date(2022, 4, 13, 17, 59, 51, 96000000, loc),
		},
		{
			name:     "nginx error",
			plugin:   func(loc *time.Location) panyl.PluginParse { return NGINXErrorLog{Location: loc} },
			source:   `2022/03/10 20:20:48 [error] 62#62: *70 invalid URL prefix in ""`,
			expected: time.Date(2022, 3, 10, 20, 20, 48, 0, loc),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, l := range []*time.Location{nil, loc} {
				ctx := context.Background()
				item := panyl.InitItem()

				ok, err := tc.plugin(l).ExtractParse(ctx, panyl.ItemLines{&panyl.Item{Line: tc.source}}, item)
				require.NoError(t, err)
				require.True(t, ok)

				expected := tc.expected
				if l == nil {
					// default is UTC
					expected = time.Date(expected.Year(), expected.Month(), expected.Day(), expected.Hour(),
						expected.Minute(), expected.Second(), expected.Nanosecond(), time.UTC)
				}
				ts, _ := item.Metadata[panyl.MetadataTimestamp].(time.Time)
				assert.True(t, expected.Equal(ts), "expected %s, got %s", expected, ts)
			}
		})
	}
}
