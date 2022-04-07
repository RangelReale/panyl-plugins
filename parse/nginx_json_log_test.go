package parse

import (
	"github.com/RangelReale/panyl"
	"github.com/RangelReale/panyl/plugins/structure"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
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
			level:   panyl.MetadataLevel_ERROR,
			message: "GET localhost:5000/graphql/query [status:500]",
		},
		{
			source:  `{"body_bytes_sent":"0","bytes_sent":"663","host":"localhost","http_host":"localhost:5000","http_request_length":"2901","http_request_method":"GET","http_request_path":"/environments.json","http_response_size":"0","http_response_time_s":"0.002","http_status_code":"304","http_user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36","message":"nginx request to /environments.json completed in 0.002 seconds","nginx_time":"1649343527.835","now":"2022-04-07T14:58:47+00:00","proxy_host":"frontend-ssr:3500","remote_addr":"172.18.0.1","request_header_size":2901,"request_length":"2901","request_method":"GET","request_time":"0.002","response_header_size":663,"status":"304","upstream_addr":"172.18.0.30:3500","upstream_bytes_received":"658","upstream_bytes_sent":"2899","upstream_connect_time":"0.000","upstream_response_time":"0.002","upstream_status":"304","uri":"/environments.json"}`,
			level:   panyl.MetadataLevel_WARNING,
			message: "GET localhost:5000/environments.json [status:304] -> upstream 172.18.0.30:3500 [status:304]",
		},
	}

	JSON := &structure.JSON{}

	for _, tc := range tests {
		result := panyl.InitProcess()
		ok, err := JSON.ExtractStructure(panyl.ProcessLines{&panyl.Process{Line: tc.source}}, result)
		assert.NoError(t, err)
		assert.True(t, ok)

		p := &NGINXJsonLog{}
		ok, err = p.ParseFormat(result)
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NotZero(t, result.Metadata[panyl.Metadata_Timestamp])
		assert.Equal(t, tc.level, result.Metadata.StringValue(panyl.Metadata_Level))
		assert.True(t, strings.HasPrefix(result.Metadata.StringValue(panyl.Metadata_Message), tc.message))
	}
}
