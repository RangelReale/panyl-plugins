package parse

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/RangelReale/panyl"
)

var _ panyl.PluginParseFormat = (*NGINXJsonLog)(nil)

const NGINXJsonLogFormat = "cb_go_json_log"

// example: "{"body_bytes_sent":"88930","bytes_sent":"89096","host":"localhost","http_host":"localhost:5000","http_request_length":"3331","http_request_method":"GET","http_request_path":"/graphql/query","http_response_size":"88930","http_response_time_s":"0.000","http_status_code":"500","http_user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36","message":"nginx request to /graphql/query completed in 0.000 seconds","nginx_time":"1649343528.035","now":"2022-04-07T14:58:48+00:00","pagekey":"homepage","platform":"web","remote_addr":"172.18.0.1","request_header_size":3331,"request_length":"3331","request_method":"GET","request_time":"0.000","response_header_size":166,"sent_http_content_type":"text/html","status":"500","uri":"/graphql/query"}"

type NGINXJsonLog struct {
}

var (
	nginxTimestampFormat = "2006-01-02T15:04:05Z07:00"
)

func (C NGINXJsonLog) ParseFormat(ctx context.Context, result *panyl.Process) (bool, error) {
	if result.Metadata.StringValue(panyl.MetadataStructure) == panyl.MetadataStructureJSON {
		if result.Data.HasValue("http_request_path") && result.Data.HasValue("http_status_code") &&
			result.Data.HasValue("nginx_time") && result.Data.HasValue("now") {

			ts, err := time.Parse(nginxTimestampFormat, result.Data.StringValue("now"))
			if err == nil {
				result.Metadata[panyl.MetadataTimestamp] = ts
			}

			level := panyl.MetadataLevelINFO
			if hsc := result.Data.StringValue("http_status_code"); hsc != "" {
				hscn, err := strconv.ParseInt(hsc, 10, 32)
				if err == nil {
					if hscn >= 400 {
						if hscn >= 500 {
							level = panyl.MetadataLevelERROR
						} else {
							level = panyl.MetadataLevelWARNING
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
			if result.Data.HasValue("proxy_host") {
				message = fmt.Sprintf("%s {proxy host:%s}", message, result.Data.StringValue("proxy_host"))
			}

			if logmessage := result.Data.StringValue("message"); message != "" {
				message = fmt.Sprintf("%s -- %s", message, logmessage)
			}

			result.Metadata[panyl.MetadataMessage] = message
			result.Metadata[panyl.MetadataLevel] = level
			result.Metadata[panyl.MetadataFormat] = NGINXJsonLogFormat
			return true, nil
		}
	}
	return false, nil
}

func (C NGINXJsonLog) IsPanylPlugin() {}
