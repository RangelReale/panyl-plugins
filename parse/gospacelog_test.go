package parse

import (
	"context"
	"fmt"
	"testing"

	"github.com/RangelReale/panyl/v2"
	"github.com/stretchr/testify/assert"
)

func TestGoSpaceLogT1(t *testing.T) {
	p := GoSpaceLog{}
	x := p.splitFields(`level=info ts=2024-12-18T14:55:27.787558447Z caller=poller.go:136 msg="blocklist poll complete" seconds=0.128523552`)
	fmt.Println(x)
}

func TestGoSpaceLog(t *testing.T) {
	type test struct {
		source   string
		level    string
		category string
		message  string
	}

	tests := []test{
		{
			source:   `level=info ts=2024-12-18T14:55:27.787558447Z caller=poller.go:136 msg="blocklist poll complete" seconds=0.128523552`,
			level:    panyl.MetadataLevelINFO,
			category: "datadog-go/tracer",
			message:  "Datadog Tracer v1.28.0 ERROR: lost 2 traces",
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
		assert.Equal(t, tc.category, item.Metadata.StringValue(panyl.MetadataCategory))
		assert.Equal(t, tc.message, item.Metadata.StringValue(panyl.MetadataMessage))
	}
}
