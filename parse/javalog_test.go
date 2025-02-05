package parse

import (
	"context"
	"testing"

	"github.com/RangelReale/panyl/v2"
	"github.com/stretchr/testify/assert"
)

func TestJavaLog(t *testing.T) {
	type test struct {
		source  string
		level   string
		message string
	}

	tests := []test{
		{
			source:  `2025-02-05 19:44:48,878 INFO Processing ruok command from /127.0.0.1:48544 (org.apache.zookeeper.server.NettyServerCnxn) [nioEventLoopGroup-4-3]`,
			level:   panyl.MetadataLevelINFO,
			message: "Processing ruok command from /127.0.0.1:48544 (org.apache.zookeeper.server.NettyServerCnxn) [nioEventLoopGroup-4-3]",
		},
	}

	for _, tc := range tests {
		ctx := context.Background()

		item := panyl.InitItem()

		p := JavaLog{}
		ok, err := p.ExtractParse(ctx, panyl.ItemLines{&panyl.Item{Line: tc.source}}, item)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, item.Metadata[panyl.MetadataTimestamp])
		assert.Equal(t, tc.level, item.Metadata.StringValue(panyl.MetadataLevel))
		assert.Equal(t, tc.message, item.Metadata.StringValue(panyl.MetadataMessage))
	}
}
