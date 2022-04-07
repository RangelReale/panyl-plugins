package parse

import (
	"fmt"
	"github.com/RangelReale/panyl"
	"strconv"
	"time"
)

var _ panyl.PluginParseFormat = (*NGINXJsonLog)(nil)

const NGINXJsonLogFormat = "cb_go_json_log"

// example: "{"body_bytes_sent":"88930","bytes_sent":"89096","host":"localhost","http_host":"localhost:5000","http_request_length":"3331","http_request_method":"GET","http_request_path":"/graphql/query","http_response_size":"88930","http_response_time_s":"0.000","http_status_code":"500","http_user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36","message":"nginx request to /graphql/query completed in 0.000 seconds","nginx_time":"1649343528.035","now":"2022-04-07T14:58:48+00:00","pagekey":"homepage","platform":"web","remote_addr":"172.18.0.1","request_header_size":3331,"request_length":"3331","request_method":"GET","request_time":"0.000","response_header_size":166,"sent_http_content_type":"text/html","status":"500","uri":"/graphql/query"}"

type NGINXJsonLog struct {
}

var (
	nginxTimestampFormat = "2006-01-02T15:04:05Z07:00"
)

func (C NGINXJsonLog) ParseFormat(result *panyl.Process) (bool, error) {
	if result.Metadata.StringValue(panyl.Metadata_Structure) == panyl.MetadataStructure_JSON {
		if result.Data.HasValue("http_request_path") && result.Data.HasValue("http_status_code") &&
			result.Data.HasValue("nginx_time") && result.Data.HasValue("now") {

			ts, err := time.Parse(nginxTimestampFormat, result.Data.StringValue("now"))
			if err == nil {
				result.Metadata[panyl.Metadata_Timestamp] = ts
			}

			level := panyl.MetadataLevel_DEBUG
			if hsc := result.Data.StringValue("http_status_code"); hsc != "" {
				hscn, err := strconv.ParseInt(hsc, 10, 32)
				if err == nil {
					if hscn >= 300 {
						if hscn >= 500 {
							level = panyl.MetadataLevel_ERROR
						} else {
							level = panyl.MetadataLevel_WARNING
						}
					}
				}
			}

			host := result.Data.StringValue("host")
			if hhost := result.Data.StringValue("http_host"); hhost != "" {
				host = hhost
			}

			message := fmt.Sprintf("%s %s%s [status:%s]",
				result.Data.StringValue("request_method"),
				host,
				result.Data.StringValue("uri"),
				result.Data.StringValue("status"),
			)
			if result.Data.HasValue("upstream_addr") {
				message = fmt.Sprintf("%s -> upstream %s [status:%s]", message,
					result.Data.StringValue("upstream_addr"),
					result.Data.StringValue("upstream_status"),
				)
			}

			if logmessage := result.Data.StringValue("message"); message != "" {
				message = fmt.Sprintf("%s -- %s", message, logmessage)
			}

			result.Metadata[panyl.Metadata_Message] = message
			result.Metadata[panyl.Metadata_Level] = level
			result.Metadata[panyl.Metadata_Format] = NGINXJsonLogFormat
			return true, nil
		}
	}
	return false, nil
}

func (C NGINXJsonLog) IsPanylPlugin() {}
