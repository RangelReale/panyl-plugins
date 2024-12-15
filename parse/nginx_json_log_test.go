package parse

import (
	"context"
	"strings"
	"testing"

	"github.com/RangelReale/panyl/v2"
	"github.com/RangelReale/panyl/v2/plugins/structure"
	"github.com/stretchr/testify/assert"
)

func TestNGINXJSONLog(t *testing.T) {
	type test struct {
		source  string
		level   string
		message string
	}

	tests := []test{
		{
			source:  `{"body_bytes_sent":"88930","bytes_sent":"89096","host":"localhost","http_host":"localhost:5000","http_request_length":"3331","http_request_method":"GET","http_request_path":"/graphql/query","http_response_size":"88930","http_response_time_s":"0.000","http_status_code":"500","http_user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36","message":"nginx request to /graphql/query completed in 0.000 seconds","nginx_time":"1649343528.035","now":"2022-04-07T14:58:48+00:00","pagekey":"homepage","platform":"web","remote_addr":"172.18.0.1","request_header_size":3331,"request_length":"3331","request_method":"GET","request_time":"0.000","response_header_size":166,"sent_http_content_type":"text/html","status":"500","uri":"/graphql/query"}`,
			level:   panyl.MetadataLevelERROR,
			message: "GET localhost:5000/graphql/query [status:500]",
		},
		{
			source:  `{"body_bytes_sent":"0","bytes_sent":"663","host":"localhost","http_host":"localhost:5000","http_request_length":"2901","http_request_method":"GET","http_request_path":"/environments.json","http_response_size":"0","http_response_time_s":"0.002","http_status_code":"404","http_user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36","message":"nginx request to /environments.json completed in 0.002 seconds","nginx_time":"1649343527.835","now":"2022-04-07T14:58:47+00:00","proxy_host":"frontend-ssr:3500","remote_addr":"172.18.0.1","request_header_size":2901,"request_length":"2901","request_method":"GET","request_time":"0.002","response_header_size":663,"status":"404","upstream_addr":"172.18.0.30:3500","upstream_bytes_received":"658","upstream_bytes_sent":"2899","upstream_connect_time":"0.000","upstream_response_time":"0.002","upstream_status":"404","uri":"/environments.json"}`,
			level:   panyl.MetadataLevelWARNING,
			message: "GET localhost:5000/environments.json [status:404] -> upstream 172.18.0.30:3500 [status:404]",
		},
		{
			source:  `{"host":"localhost","request_method":"GET","status":"200","bytes_sent":"129299","message":"nginx request to /dashboard completed in 0.000 seconds","http_host":"localhost:5000","response_header_size":443,"nginx_time":"1649343527.265","http_response_size":"128856","uri":"/dashboard","now":"2022-04-07T14:58:47+00:00","request_length":"2979","proxy_host":"s3.amazonaws.com","http_response_time_s":"0.000","http_user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36","http_request_method":"GET","http_request_path":"/dashboard","request_time":"0.000","http_request_length":"2979","request_header_size":2979,"remote_addr":"172.18.0.1","sent_http_content_type":"text/html; charset=utf-8","upstream_cache_status":"HIT","http_status_code":"200","body_bytes_sent":"128856"}`,
			level:   panyl.MetadataLevelINFO,
			message: "GET localhost:5000/dashboard [status:200]",
		},
	}

	JSON := &structure.JSON{}

	for _, tc := range tests {
		ctx := context.Background()
		item := panyl.InitItem()
		ok, err := JSON.ExtractStructure(ctx, panyl.ItemLines{&panyl.Item{Line: tc.source}}, item)
		assert.NoError(t, err)
		assert.True(t, ok)

		p := &NGINXJsonLog{}
		ok, err = p.ParseFormat(ctx, item)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, item.Metadata[panyl.MetadataTimestamp])
		assert.Equal(t, tc.level, item.Metadata.StringValue(panyl.MetadataLevel))
		assert.True(t, strings.HasPrefix(item.Metadata.StringValue(panyl.MetadataMessage), tc.message))
	}
}
