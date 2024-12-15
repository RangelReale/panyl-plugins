package metadata

import (
	"context"
	"regexp"
	"strings"

	"github.com/RangelReale/panyl/v2"
)

// DockerCompose extracts application name from the line by the docker compose format, which is an application
// name followed by | at the beginning of the line.
// It also signals a sequence break on lines of different applications.
// If ApplicationWhitelist is not nil, only applications on this list will be considered.
type DockerCompose struct {
	OnlyIfAnsiEscape     bool
	ApplicationWhitelist []string
}

var _ panyl.PluginMetadata = (*DockerCompose)(nil)
var _ panyl.PluginSequence = (*DockerCompose)(nil)

// example: "application    |"

var dockerPrefixRE = regexp.MustCompile(`^(\w|[-])+\s+\|`)

func (m *DockerCompose) ExtractMetadata(ctx context.Context, item *panyl.Item) (bool, error) {
	matches := dockerPrefixRE.FindStringSubmatchIndex(item.Line)
	if matches == nil {
		return false, nil
	}

	if m.OnlyIfAnsiEscape && !item.Metadata.ListValueContains(panyl.MetadataClean, panyl.MetadataCleanAnsiEscape) {
		return false, nil
	}

	application := strings.TrimSpace(item.Line[matches[0] : matches[1]-1])
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

	if len(item.Line) > matches[1] {
		item.Line = item.Line[matches[1]+1:]
		return true, nil
	}
	item.Line = ""
	return true, nil
}

func (m *DockerCompose) BlockSequence(ctx context.Context, lastp, item *panyl.Item) bool {
	// block sequence if application changed
	return lastp.Metadata.StringValue(panyl.MetadataApplication) != item.Metadata.StringValue(panyl.MetadataApplication)
}

func (m DockerCompose) IsPanylPlugin() {}
