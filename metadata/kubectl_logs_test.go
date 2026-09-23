package metadata

import (
	"context"
	"testing"

	"github.com/RangelReale/panyl/v2"
	"github.com/stretchr/testify/assert"
)

func TestKubeCtlLogs(t *testing.T) {
	tests := []struct {
		source      string
		extractName bool
		application string
		line        string
	}{
		{source: "[pod/events-worker-7c9b7bdc55-f7sgc/worker] log text", application: "pod/events-worker-7c9b7bdc55-f7sgc/worker", line: "log text"},
		{source: "[pod/events-worker-7c9b7bdc55-f7sgc/worker] log text", extractName: true, application: "events-worker/worker", line: "log text"},
		{source: "[pod/init-database-fm44h/init] log text", extractName: true, application: "init-database/init", line: "log text"},
		{source: "[pod/api-proxy/proxy] log text", extractName: true, application: "api-proxy/proxy", line: "log text"},
		{source: "[pod/deadbeef-abcd/app] log text", extractName: true, application: "deadbeef-abcd/app", line: "log text"},
		{source: "[pod/cafebabe-f7sgc/app] log text", extractName: true, application: "cafebabe/app", line: "log text"},
		{source: "[pod/redis/redis] log text", extractName: true, application: "redis/redis", line: "log text"},
	}

	for _, tc := range tests {
		ctx := context.Background()
		item := panyl.InitItem()
		item.Line = tc.source

		ok, err := KubeCtlLogs{ExtractApplicationName: tc.extractName}.ExtractMetadata(ctx, item)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.Equal(t, tc.application, item.Metadata.StringValue(panyl.MetadataApplication), tc.source)
		assert.Equal(t, tc.line, item.Line)
	}
}
