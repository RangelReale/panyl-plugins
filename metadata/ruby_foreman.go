package metadata

import (
	"github.com/RangelReale/panyl"
	"regexp"
	"strings"
)

var _ panyl.PluginMetadata = (*RubyForeman)(nil)
var _ panyl.PluginSequence = (*RubyForeman)(nil)

// RubyForeman extracts application name from the line by the roby foreman format, which is
// a time, followed by an application
// name, followed by | at the beginning of the line.
// It also signals a sequence break on lines of different applications.
// If ApplicationWhitelist is not nil, only applications on this list will be considered.
type RubyForeman struct {
	OnlyIfAnsiEscape     bool
	ApplicationWhitelist []string
}

// example: "16:41:59 api.1         | log text"

var rubyForemanPrefixRE = regexp.MustCompile(`^(\d{2}:\d{2}:\d{2})\s([\w.]+)\s+\|(.*)$`)

func (m *RubyForeman) ExtractMetadata(result *panyl.Process) (bool, error) {
	matches := rubyForemanPrefixRE.FindStringSubmatch(result.Line)
	if matches == nil {
		return false, nil
	}

	if m.OnlyIfAnsiEscape && !result.Metadata.ListValueContains(panyl.MetadataClean, panyl.MetadataCleanAnsiEscape) {
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

	result.Metadata[panyl.MetadataApplication] = application
	if len(text) > 0 {
		result.Line = text
		return true, nil
	}
	result.Line = ""
	return true, nil
}

func (m *RubyForeman) BlockSequence(lastp, p *panyl.Process) bool {
	// block sequence if application changed
	return lastp.Metadata.StringValue(panyl.MetadataApplication) != p.Metadata.StringValue(panyl.MetadataApplication)
}

func (m RubyForeman) IsPanylPlugin() {}
