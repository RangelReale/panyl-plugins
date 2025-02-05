package sequence

import (
	"context"

	"github.com/RangelReale/panyl/v2"
)

// RawLine returns all lines.
type RawLine struct {
	SourceAsCategory bool
}

var _ panyl.PluginSequence = RawLine{}

func (m RawLine) BlockSequence(ctx context.Context, lastp, item *panyl.Item) bool {
	return true
}

func (m RawLine) IsPanylPlugin() {}
