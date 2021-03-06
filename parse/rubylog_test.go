package parse

import (
	"github.com/RangelReale/panyl"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRubyLog(t *testing.T) {
	type test struct {
		source   string
		level    string
		category string
		message  string
	}

	tests := []test{
		{
			source:   `I, [1999-03-03T02:34:24.895701 #19074]  INFO -- Main: info.`,
			level:    panyl.MetadataLevel_INFO,
			category: "",
			message:  "info.",
		},
	}

	for _, tc := range tests {
		result := panyl.InitProcess()

		p := &RubyLog{}
		ok, err := p.ExtractParse(panyl.ProcessLines{&panyl.Process{Line: tc.source}}, result)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, result.Metadata[panyl.Metadata_Timestamp])
		assert.Equal(t, tc.level, result.Metadata.StringValue(panyl.Metadata_Level))
		assert.Equal(t, tc.category, result.Metadata.StringValue(panyl.Metadata_Category))
		assert.Equal(t, tc.message, result.Metadata.StringValue(panyl.Metadata_Message))
	}
}
