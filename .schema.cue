import (
	"strings"
	"time"
)

#Metadata: {
	title:        _#Title @error(title is required)
	type:         _#Type  @error(type is one of either: concept, guide, or reference)
	lastReviewed: _#Date  @error(document must be reviewed on, or more recently than, Jan 1st 2025)
}

_#Date: time.Format(time.RFC3339Date) & >="2025-01-01"

_#Title: string & strings.MinRunes(1)

_#Type: "concept" | "guide" | "reference"
