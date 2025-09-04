package parse

import (
	"context"
	"testing"

	"github.com/RangelReale/panyl/v2"
	"github.com/RangelReale/panyl/v2/plugins/structure"
	"github.com/stretchr/testify/assert"
)

func TestKubeEventJSONLog(t *testing.T) {
	type test struct {
		source  string
		level   string
		message string
	}

	tests := []test{
		{
			source:  `{"apiVersion":"v1","count":1,"eventTime":null,"firstTimestamp":"2025-08-25T11:47:00Z","involvedObject":{"apiVersion":"batch/v1","kind":"Job","name":"server-wake-scheduler-29268707","namespace":"company-apps","resourceVersion":"1123395957","uid":"bda6fdf0-7b63-4bb3-ba8d-7a2bc3e63181"},"kind":"Event","lastTimestamp":"2025-08-25T11:47:00Z","message":"Created pod: server-wake-scheduler-29268707-qwnzj","metadata":{"creationTimestamp":"2025-08-25T11:47:00Z","name":"server-wake-scheduler-29268707.185f00096aa7f57d","namespace":"company-apps","resourceVersion":"15806623","uid":"bcb0a277-0313-443f-9a1d-8f006a3bebe7"},"reason":"SuccessfulCreate","reportingComponent":"job-controller","reportingInstance":"","source":{"component":"job-controller"},"type":"Normal"}`,
			level:   panyl.MetadataLevelINFO,
			message: "(Job) [company-apps:server-wake-scheduler-29268707] Created pod: server-wake-scheduler-29268707-qwnzj [reason:SuccessfulCreate]",
		},
	}

	JSON := structure.JSON{}

	for _, tc := range tests {
		ctx := context.Background()
		item := panyl.InitItem()
		ok, err := JSON.ExtractStructure(ctx, panyl.ItemLines{&panyl.Item{Line: tc.source}}, item)
		assert.NoError(t, err)
		assert.True(t, ok)

		p := KubeEventJsonLog{}
		ok, err = p.ParseFormat(ctx, item)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, item.Metadata[panyl.MetadataTimestamp])
		assert.Equal(t, tc.level, item.Metadata.StringValue(panyl.MetadataLevel))
		assert.Equal(t, tc.message, item.Metadata.StringValue(panyl.MetadataMessage))
	}
}
