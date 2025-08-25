package parse

import (
	"context"
	"fmt"
	"time"

	"github.com/RangelReale/panyl/v2"
)

const KubeEventJsonLogFormat = "cb_kube_json_log"

// example: "{"apiVersion":"v1","count":1,"eventTime":null,"firstTimestamp":"2025-08-25T11:47:00Z","involvedObject":{"apiVersion":"batch/v1","kind":"Job","name":"server-wake-scheduler-29268707","namespace":"company-apps","resourceVersion":"1123395957","uid":"bda6fdf0-7b63-4bb3-ba8d-7a2bc3e63181"},"kind":"Event","lastTimestamp":"2025-08-25T11:47:00Z","message":"Created pod: server-wake-scheduler-29268707-qwnzj","metadata":{"creationTimestamp":"2025-08-25T11:47:00Z","name":"server-wake-scheduler-29268707.185f00096aa7f57d","namespace":"company-apps","resourceVersion":"15806623","uid":"bcb0a277-0313-443f-9a1d-8f006a3bebe7"},"reason":"SuccessfulCreate","reportingComponent":"job-controller","reportingInstance":"","source":{"component":"job-controller"},"type":"Normal"}"

type KubeEventJsonLog struct {
}

var _ panyl.PluginParseFormat = KubeEventJsonLog{}

var (
	// kubeEventTimestampFormat = "2006-01-02T15:04:05Z07:00"
	kubeEventTimestampFormat = "2006-01-02T15:04:05Z"
)

func (m KubeEventJsonLog) ParseFormat(ctx context.Context, item *panyl.Item) (bool, error) {
	if item.Metadata.StringValue(panyl.MetadataStructure) == panyl.MetadataStructureJSON {
		if item.Data.HasValues("type", "apiVersion", "involvedObject", "kind", "lastTimestamp", "message") &&
			item.Data.StringValue("kind") == "Event" {
			ts, err := time.Parse(kubeEventTimestampFormat, item.Data.StringValue("lastTimestamp"))
			if err == nil {
				item.Metadata[panyl.MetadataTimestamp] = ts
			}

			level := panyl.MetadataLevelINFO
			if hsc := item.Data.StringValue("type"); hsc != "" {
				switch hsc {
				case "Normal":
					level = panyl.MetadataLevelINFO
				case "Warning":
					level = panyl.MetadataLevelWARNING
				case "Error":
					level = panyl.MetadataLevelERROR
				default:
					level = panyl.MetadataLevelERROR
				}
			}

			var message string

			if logmessage := item.Data.StringValue("message"); logmessage != "" {
				message = logmessage
			}

			if involvedObject := item.Data.MapValue("involvedObject"); involvedObject != nil {
				objectName := fmt.Sprintf("[%s:%s](%s/%s)",
					involvedObject.StringValue("apiVersion"), involvedObject.StringValue("kind"),
					involvedObject.StringValue("namespace"), involvedObject.StringValue("name"))
				message = fmt.Sprintf("%s %s", objectName, message)
			}
			if rs := item.Data.StringValue("reason"); rs != "" {
				message = fmt.Sprintf("%s [reason:%s]", message, rs)
			}

			item.Metadata[panyl.MetadataMessage] = message
			item.Metadata[panyl.MetadataLevel] = level
			item.Metadata[panyl.MetadataFormat] = KubeEventJsonLogFormat
			return true, nil
		}
	}
	return false, nil
}

func (m KubeEventJsonLog) IsPanylPlugin() {}
