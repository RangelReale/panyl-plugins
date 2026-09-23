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
		{
			source:  `2025-02-05 19:44:48,878 WARNING Connection lost`,
			level:   panyl.MetadataLevelWARNING,
			message: "Connection lost",
		},
		{
			source:  `2025-02-05 19:44:48,878 FATAL Out of memory`,
			level:   panyl.MetadataLevelERROR,
			message: "Out of memory",
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

func TestJavaLogNoMatch(t *testing.T) {
	sources := []string{
		// Postgres line: "UTC" must not be accepted as a level
		`2022-04-05 14:29:07.500 UTC [73] ERROR:  relation "users" does not exist at character 36`,
	}

	for _, source := range sources {
		ctx := context.Background()
		item := panyl.InitItem()

		p := JavaLog{}
		ok, err := p.ExtractParse(ctx, panyl.ItemLines{&panyl.Item{Line: source}}, item)
		assert.NoError(t, err)
		assert.False(t, ok, source)
	}
}
