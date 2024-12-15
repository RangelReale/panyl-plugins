package metadata

import (
	"context"
	"regexp"
	"strings"

	"github.com/RangelReale/panyl/v2"
)

// RubyForeman extracts application name from the line by the roby foreman format, which is
// a time, followed by an application
// name, followed by | at the beginning of the line.
// It also signals a sequence break on lines of different applications.
// If ApplicationWhitelist is not nil, only applications on this list will be considered.
type RubyForeman struct {
	OnlyIfAnsiEscape     bool
	ApplicationWhitelist []string
}

var _ panyl.PluginMetadata = (*RubyForeman)(nil)
var _ panyl.PluginSequence = (*RubyForeman)(nil)

// example: "16:41:59 api.1         | log text"

var rubyForemanPrefixRE = regexp.MustCompile(`^(\d{2}:\d{2}:\d{2})\s([\w.]+)\s+\|(.*)$`)

func (m *RubyForeman) ExtractMetadata(ctx context.Context, item *panyl.Item) (bool, error) {
	matches := rubyForemanPrefixRE.FindStringSubmatch(item.Line)
	if matches == nil {
		return false, nil
	}

	if m.OnlyIfAnsiEscape && !item.Metadata.ListValueContains(panyl.MetadataClean, panyl.MetadataCleanAnsiEscape) {
		return false, nil
	}

	// time := matches[1]
	application := strings.TrimSpace(matches[2])
	text := strings.TrimSpace(matches[3])

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
	if len(text) > 0 {
		item.Line = text
		return true, nil
	}
	item.Line = ""
	return true, nil
}

func (m *RubyForeman) BlockSequence(ctx context.Context, lastp, item *panyl.Item) bool {
	// block sequence if application changed
	return lastp.Metadata.StringValue(panyl.MetadataApplication) != item.Metadata.StringValue(panyl.MetadataApplication)
}

func (m RubyForeman) IsPanylPlugin() {}
