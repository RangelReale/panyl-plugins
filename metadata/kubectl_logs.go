package metadata

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/RangelReale/panyl/v2"
)

const (
	KubeCtlLogsApplicationSource = "kubectllogs"
)

// KubeCtlLogs extracts application name from the line by the "kubectl logs --prefix format", which is pod name
// between [] at the beginning of the line.
// It also signals a sequence break on lines of different applications.
// If ApplicationWhitelist is not nil, only applications on this list will be considered.
type KubeCtlLogs struct {
	OnlyIfAnsiEscape       bool
	ExtractApplicationName bool

	ApplicationWhitelist []string
}

var _ panyl.PluginMetadata = DockerCompose{}
var _ panyl.PluginSequence = DockerCompose{}

// example: "[pod/deployment-id1-id2/container] "

// pod names:
// events-worker-7c9b7bdc55-f7sgc
// notification-585f6b94b8-t4jjc
// report-generator-749ccf648d-gd9j7
// static-content-8445d7f556-wbxvp
// init-database-fm44h

var (
	kubeCtlLogsPrefixRE = regexp.MustCompile(`^\[([^\]]+)\] `)
	kubeCtlLogsHexRE    = regexp.MustCompile(`[^0-9A-Fa-f]`)
)

func (m KubeCtlLogs) ExtractMetadata(ctx context.Context, item *panyl.Item) (bool, error) {
	matches := kubeCtlLogsPrefixRE.FindStringSubmatch(item.Line)
	if matches == nil {
		return false, nil
	}

	if m.OnlyIfAnsiEscape && !item.Metadata.ListValueContains(panyl.MetadataClean, panyl.MetadataCleanAnsiEscape) {
		return false, nil
	}

	application := strings.TrimSpace(matches[1])
	appsep := strings.Split(application, "/")
	if m.ExtractApplicationName && len(appsep) == 3 {
		application = m.parsePodName(appsep[1], appsep[2])
		// matches := kubeCtlLogsDeploymentRE.FindStringSubmatch(appsep[1])
		// if matches != nil {
		// 	application = fmt.Sprintf("%s/%s", matches[2], appsep[2])
		// }
	}

	if len(m.ApplicationWhitelist) > 0 {
		found := false
		for _, app := range m.ApplicationWhitelist {
			if app == application {
				found = true
				break
			}
		}
		if !found {
			return false, nil
		}
	}

	item.Metadata[panyl.MetadataApplication] = application
	item.Metadata[panyl.MetadataApplicationSource] = KubeCtlLogsApplicationSource

	if len(item.Line) > len(matches[1]) {
		item.Line = item.Line[len(matches[1])+3:]
		return true, nil
	}
	item.Line = ""
	return true, nil
}

func (m KubeCtlLogs) BlockSequence(ctx context.Context, lastp, item *panyl.Item) bool {
	// block sequence if application changed
	return lastp.Metadata.StringValue(panyl.MetadataApplication) != item.Metadata.StringValue(panyl.MetadataApplication)
}

func (m KubeCtlLogs) IsPanylPlugin() {}

func (m KubeCtlLogs) parsePodName(name string, containerName string) string {
	ns := strings.Split(name, "-")
	if len(ns) < 2 {
		return name
	}

	// last item may have 4 or 5 chars
	lastLen := len(ns[len(ns)-1])
	if lastLen == 4 || lastLen == 5 {
		ns = slices.Delete(ns, len(ns)-1, len(ns))
	}

	// last item removed only if 8 or 10 chars, and a hex string
	lastPrev := ns[len(ns)-1]
	if len(lastPrev) >= 8 && len(lastPrev) <= 10 {
		// "true" means that non-hex charts exist
		if !kubeCtlLogsHexRE.MatchString(lastPrev) {
			// s is a valid hex string
			ns = slices.Delete(ns, len(ns)-1, len(ns))
		}
	}

	return fmt.Sprintf("%s/%s", strings.Join(ns, "-"), containerName)
}
