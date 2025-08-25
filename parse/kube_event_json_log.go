package parse

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/RangelReale/panyl/v2"
)

const KubeEventJsonLogFormat = "cb_kube_json_log"

// example: "{"apiVersion":"v1","count":1,"eventTime":null,"firstTimestamp":"2025-08-25T11:47:00Z","involvedObject":{"apiVersion":"batch/v1","kind":"Job","name":"server-wake-scheduler-29268707","namespace":"company-apps","resourceVersion":"1123395957","uid":"bda6fdf0-7b63-4bb3-ba8d-7a2bc3e63181"},"kind":"Event","lastTimestamp":"2025-08-25T11:47:00Z","message":"Created pod: server-wake-scheduler-29268707-qwnzj","metadata":{"creationTimestamp":"2025-08-25T11:47:00Z","name":"server-wake-scheduler-29268707.185f00096aa7f57d","namespace":"company-apps","resourceVersion":"15806623","uid":"bcb0a277-0313-443f-9a1d-8f006a3bebe7"},"reason":"SuccessfulCreate","reportingComponent":"job-controller","reportingInstance":"","source":{"component":"job-controller"},"type":"Normal"}"

type KubeEventJsonLog struct {
	NamespaceAsCategory bool
	SkipObjects         []string
	SkipReasons         []string
}

var _ panyl.PluginParseFormat = KubeEventJsonLog{}

var (
	// kubeEventTimestampFormat = "2006-01-02T15:04:05Z07:00"
	kubeEventTimestampFormat = "2006-01-02T15:04:05Z"
)

func (m KubeEventJsonLog) ParseFormat(ctx context.Context, item *panyl.Item) (bool, error) {
	if item.Metadata.StringValue(panyl.MetadataStructure) == panyl.MetadataStructureJSON {
		if item.Data.HasValues("type", "apiVersion", "involvedObject", "kind") &&
			item.Data.StringValue("kind") == "Event" {
			if ts, err := time.Parse(kubeEventTimestampFormat, item.Data.StringValue("eventTime")); err == nil {
				item.Metadata[panyl.MetadataTimestamp] = ts
			} else if ts, err := time.Parse(kubeEventTimestampFormat, item.Data.StringValue("lastTimestamp")); err == nil {
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
			var category string

			if logmessage := item.Data.StringValue("message"); logmessage != "" {
				message = logmessage
			}

			if ct := item.Data.IntValue("count"); ct > 1 {
				ctmsg := fmt.Sprintf("{ct:%d}", ct)
				if message != "" {
					message = fmt.Sprintf("%s %s", ctmsg, message)
				} else {
					message = ctmsg
				}
			}

			isBatch := false

			if involvedObject := item.Data.MapValue("involvedObject"); involvedObject != nil {
				objectType := involvedObject.StringValue("kind")
				objectName := involvedObject.StringValue("name")
				namespacedObjectType := objectType
				if objectType != "" {
					if apiVersion := involvedObject.StringValue("apiVersion"); apiVersion != "" {
						namespacedObjectType = fmt.Sprintf("%s:%s", apiVersion, objectType)
						if !slices.Contains([]string{"v1", "apps/v1", "batch/v1"}, apiVersion) {
							objectType = fmt.Sprintf("%s:%s", apiVersion, objectType)
						}
						if apiVersion == "batch/v1" {
							isBatch = true
						}
					}
				}
				if ns := involvedObject.StringValue("namespace"); ns != "" {
					if m.NamespaceAsCategory {
						category = ns
					} else if objectName != "" {
						objectName = fmt.Sprintf("%s:%s", ns, objectName)
					}
				}
				if namespacedObjectType != "" && slices.Contains(m.SkipObjects, namespacedObjectType) {
					item.Metadata[panyl.MetadataSkip] = true
				}

				var objectFullname string
				if objectType != "" {
					objectFullname += fmt.Sprintf("(%s)", objectType)
				}
				if objectName != "" {
					objectFullname += fmt.Sprintf(" [%s]", objectName)
				}

				if objectFullname != "" {
					if message != "" {
						message = fmt.Sprintf("%s %s", objectFullname, message)
					} else {
						message = objectFullname
					}
				}
			}
			if rs := item.Data.StringValue("reason"); rs != "" {
				if message != "" {
					message = fmt.Sprintf("%s [reason:%s]", message, rs)
				} else {
					message = rs
				}
				if slices.Contains(m.SkipReasons, rs) {
					item.Metadata[panyl.MetadataSkip] = true
				}
			}

			if isBatch {
				category = "batch"
			}

			item.Metadata[panyl.MetadataMessage] = message
			item.Metadata[panyl.MetadataLevel] = level
			item.Metadata[panyl.MetadataFormat] = KubeEventJsonLogFormat
			if category != "" {
				item.Metadata[panyl.MetadataCategory] = category
			}
			return true, nil
		}
	}
	return false, nil
}

func (m KubeEventJsonLog) IsPanylPlugin() {}
