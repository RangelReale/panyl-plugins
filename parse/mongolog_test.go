package parse

import (
	"context"
	"testing"

	"github.com/RangelReale/panyl/v2"
	"github.com/stretchr/testify/assert"
)

func TestMongoLog(t *testing.T) {
	type test struct {
		source   string
		level    string
		category string
		message  string
	}

	tests := []test{
		{
			source:   `2022-03-10T20:17:25.039+0000 I  NETWORK  [conn15] end connection 172.19.0.1:49812 (7 connections now open)`,
			level:    panyl.MetadataLevelINFO,
			category: "NETWORK",
			message:  "end connection 172.19.0.1:49812 (7 connections now open)",
		},
	}

	for _, tc := range tests {
		ctx := context.Background()
		result := panyl.InitProcess()

		p := &MongoLog{}
		ok, err := p.ExtractParse(ctx, panyl.ProcessLines{&panyl.Process{Line: tc.source}}, result)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, result.Metadata[panyl.MetadataTimestamp])
		assert.Equal(t, tc.level, result.Metadata.StringValue(panyl.MetadataLevel))
		assert.Equal(t, tc.category, result.Metadata.StringValue(panyl.MetadataCategory))
		assert.Equal(t, tc.message, result.Metadata.StringValue(panyl.MetadataMessage))
	}
}
