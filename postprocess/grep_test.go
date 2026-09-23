package postprocess

import (
	"context"
	"testing"

	"github.com/RangelReale/panyl/v2"
	"github.com/stretchr/testify/assert"
)

func TestGrep(t *testing.T) {
	tests := []struct {
		values   []string
		message  string
		expected []string
	}{
		{values: []string{"hello"}, message: "Hello world", expected: []string{"[grep][hello]"}},
		{values: []string{"foo"}, message: "Hello world", expected: nil},
		{values: []string{"", " ", "foo"}, message: "Hello world", expected: nil},
	}

	for _, tc := range tests {
		ctx := context.Background()
		item := panyl.InitItem()
		item.Metadata[panyl.MetadataMessage] = tc.message

		_, err := Grep{Values: tc.values}.PostProcess(ctx, item)
		assert.NoError(t, err)

		assert.Equal(t, tc.expected, item.Metadata.ListValue(panyl.MetadataExtraCategories))
	}
}
