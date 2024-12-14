package parse

import (
	"github.com/RangelReale/panyl"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoLog(t *testing.T) {
	type test struct {
		source   string
		level    string
		category string
		message  string
	}

	tests := []test{
		{
			source:   `2022-03-10T19:53:21.434Z	INFO	datadog-go/tracer.go:35	Datadog Tracer v1.28.0 ERROR: lost 2 traces`,
			level:    panyl.MetadataLevelINFO,
			category: "datadog-go/tracer",
			message:  "Datadog Tracer v1.28.0 ERROR: lost 2 traces",
		},
	}

	for _, tc := range tests {
		result := panyl.InitProcess()

		p := &GoLog{SourceAsCategory: true}
		ok, err := p.ExtractParse(panyl.ProcessLines{&panyl.Process{Line: tc.source}}, result)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, result.Metadata[panyl.MetadataTimestamp])
		assert.Equal(t, tc.level, result.Metadata.StringValue(panyl.MetadataLevel))
		assert.Equal(t, tc.category, result.Metadata.StringValue(panyl.MetadataCategory))
		assert.Equal(t, tc.message, result.Metadata.StringValue(panyl.MetadataMessage))
	}
}
