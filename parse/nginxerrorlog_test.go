package parse

import (
	"github.com/RangelReale/panyl"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNGINXErrorLog(t *testing.T) {
	type test struct {
		source   string
		level    string
		category string
		message  string
	}

	tests := []test{
		{
			source:   `2022/03/10 20:20:48 [error] 62#62: *70 invalid URL prefix in "", client: 127.0.0.1, server: , request: "GET / HTTP/1.1", host: "localhost:8080", referrer: "http://localhost:8080"`,
			level:    panyl.MetadataLevel_ERROR,
			category: "",
			message:  "invalid URL prefix in",
		},
	}

	for _, tc := range tests {
		result := panyl.InitProcess()

		p := &NGINXErrorLog{}
		ok, err := p.ExtractParse(panyl.ProcessLines{&panyl.Process{Line: tc.source}}, result)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, result.Metadata[panyl.Metadata_Timestamp])
		assert.Equal(t, tc.level, result.Metadata.StringValue(panyl.Metadata_Level))
		assert.Equal(t, tc.category, result.Metadata.StringValue(panyl.Metadata_Category))
		assert.True(t, strings.HasPrefix(result.Metadata.StringValue(panyl.Metadata_Message), tc.message),
			"message is different, expected prefix '%s' got '%s'", tc.message,
			result.Metadata.StringValue(panyl.Metadata_Message))
	}
}
