package utils

import "strings"

// ShortText collapses runs of whitespace into single spaces and truncates the
// result to at most max runes (appending "…" when truncated). max <= 0 means
// no truncation. Truncation is rune-safe: it never splits a multi-byte UTF-8
// sequence. Used for log summaries (e.g. SureSQL's "medium" log level) and any
// other place a compact one-line preview of free text is needed.
func ShortText(s string, max int) string {
	s = strings.Join(strings.Fields(s), " ")
	if max <= 0 || len(s) <= max {
		return s
	}
	// Count runes; cut only on a rune boundary.
	n := 0
	for i := range s {
		if n == max {
			return s[:i] + "…"
		}
		n++
	}
	return s
}
