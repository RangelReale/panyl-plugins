package parse

import (
	"context"
	"strings"
	"testing"

	"github.com/RangelReale/panyl/v2"
	"github.com/stretchr/testify/assert"
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
			level:    panyl.MetadataLevelERROR,
			category: "",
			message:  "invalid URL prefix in",
		},
	}

	for _, tc := range tests {
		ctx := context.Background()
		item := panyl.InitItem()

		p := &NGINXErrorLog{}
		ok, err := p.ExtractParse(ctx, panyl.ItemLines{&panyl.Item{Line: tc.source}}, item)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, item.Metadata[panyl.MetadataTimestamp])
		assert.Equal(t, tc.level, item.Metadata.StringValue(panyl.MetadataLevel))
		assert.Equal(t, tc.category, item.Metadata.StringValue(panyl.MetadataCategory))
		assert.True(t, strings.HasPrefix(item.Metadata.StringValue(panyl.MetadataMessage), tc.message),
			"message is different, expected prefix '%s' got '%s'", tc.message,
			item.Metadata.StringValue(panyl.MetadataMessage))
	}
}
