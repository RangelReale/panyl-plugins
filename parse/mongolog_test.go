package parse

import (
	"github.com/RangelReale/panyl"
	"github.com/stretchr/testify/assert"
	"testing"
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
			level:    panyl.MetadataLevel_INFO,
			category: "NETWORK",
			message:  "end connection 172.19.0.1:49812 (7 connections now open)",
		},
	}

	for _, tc := range tests {
		result := panyl.InitProcess()

		p := &MongoLog{}
		ok, err := p.ExtractParse(panyl.ProcessLines{&panyl.Process{Line: tc.source}}, result)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, result.Metadata[panyl.Metadata_Timestamp])
		assert.Equal(t, tc.level, result.Metadata.StringValue(panyl.Metadata_Level))
		assert.Equal(t, tc.category, result.Metadata.StringValue(panyl.Metadata_Category))
		assert.Equal(t, tc.message, result.Metadata.StringValue(panyl.Metadata_Message))
	}
}
