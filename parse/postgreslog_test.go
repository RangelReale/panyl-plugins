package parse

import (
	"github.com/RangelReale/panyl"
	"github.com/stretchr/testify/assert"
	"testing"
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
			level:    panyl.MetadataLevel_ERROR,
			category: "",
			message:  `relation "users" does not exist at character 36`,
		},
	}

	for _, tc := range tests {
		result := panyl.InitProcess()

		p := &PostgresLog{}
		ok, err := p.ExtractParse(panyl.ProcessLines{&panyl.Process{Line: tc.source}}, result)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, result.Metadata[panyl.Metadata_Timestamp])
		assert.Equal(t, tc.level, result.Metadata.StringValue(panyl.Metadata_Level))
		assert.Equal(t, tc.category, result.Metadata.StringValue(panyl.Metadata_Category))
		assert.Equal(t, tc.message, result.Metadata.StringValue(panyl.Metadata_Message))
	}
}
