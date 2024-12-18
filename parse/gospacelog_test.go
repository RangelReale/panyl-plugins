package parse

import (
	"context"
	"testing"

	"github.com/RangelReale/panyl/v2"
	gocmp "github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoSpaceLog(t *testing.T) {
	type test struct {
		source  string
		level   string
		message string
		data    panyl.MapValue
	}

	tests := []test{
		{
			source:  `level=info ts=2024-12-18T14:55:27.787558447Z caller=poller.go:136 msg="blocklist poll complete" seconds=0.128523552`,
			level:   panyl.MetadataLevelINFO,
			message: "blocklist poll complete",
			data: panyl.MapValue{
				"caller":  "poller.go:136",
				"msg":     "blocklist poll complete",
				"seconds": "0.128523552",
				"ts":      "2024-12-18T14:55:27.787558447Z",
				"level":   "info",
			},
		},
		{
			source:  `level=info ts=2024-12-18T14:55:27.787558447Z caller=poller.go:136 msg="blocklist \"with quotes\" poll complete" seconds=0.128523552`,
			level:   panyl.MetadataLevelINFO,
			message: `blocklist "with quotes" poll complete`,
			data: panyl.MapValue{
				"caller":  "poller.go:136",
				"msg":     `blocklist "with quotes" poll complete`,
				"seconds": "0.128523552",
				"ts":      "2024-12-18T14:55:27.787558447Z",
				"level":   "info",
			},
		},
		{
			source:  `level=info ts=2024-12-18T14:55:27.787558447Z caller=poller.go:136 msg="blocklist tab \t poll complete" seconds=0.128523552`,
			level:   panyl.MetadataLevelINFO,
			message: "blocklist tab \t poll complete",
			data: panyl.MapValue{
				"caller":  "poller.go:136",
				"msg":     "blocklist tab \t poll complete",
				"seconds": "0.128523552",
				"ts":      "2024-12-18T14:55:27.787558447Z",
				"level":   "info",
			},
		},
	}

	for _, tc := range tests {
		ctx := context.Background()

		item := panyl.InitItem()

		p := GoSpaceLog{SourceAsCategory: true}
		ok, err := p.ExtractParse(ctx, panyl.ItemLines{&panyl.Item{Line: tc.source}}, item)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, item.Metadata[panyl.MetadataTimestamp])
		assert.Equal(t, tc.level, item.Metadata.StringValue(panyl.MetadataLevel))
		assert.Equal(t, tc.message, item.Metadata.StringValue(panyl.MetadataMessage))
		requireDeepEqual(t, tc.data, item.Data)
	}
}

type helperT interface {
	Helper()
}

// assertDeepEqual uses [github.com/google/go-cmp/cmp]
// to assert two values are equal and fails the test if they are not equal.
func assertDeepEqual(t require.TestingT, expected interface{}, actual interface{}, opts ...gocmp.Option) bool {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	diff := gocmp.Diff(expected, actual, opts...)
	if diff == "" {
		return true
	}

	return assert.Fail(t, diff)
}

// requireDeepEqual uses [github.com/google/go-cmp/cmp]
// to assert two values are equal and fails the test if they are not equal.
func requireDeepEqual(t require.TestingT, expected interface{}, actual interface{}, opts ...gocmp.Option) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}
	if assertDeepEqual(t, expected, actual, opts...) {
		return
	}
	t.FailNow()
}
