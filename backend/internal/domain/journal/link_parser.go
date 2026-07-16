package journal

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

const previewMaxLength = 180

var (
	customMoodRecordLinkPattern = regexp.MustCompile(`\[\[mood:(\d+)(?:\|([^\]]+))?\]\]`)
	legacyMoodRecordLinkPattern = regexp.MustCompile(`(?:https?://[^\s)]+)?/mood/records/(\d+)(?:[?#][^\s)]*)?`)
	fencedCodeBlockPattern      = regexp.MustCompile("(?s)```.*?```|~~~.*?~~~")
	inlineCodePattern           = regexp.MustCompile("`[^`\n]*`")
	markdownLinkPattern         = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	markdownHeadingPattern      = regexp.MustCompile(`(?m)^\s{0,3}#{1,6}\s*`)
	markdownQuotePattern        = regexp.MustCompile(`(?m)^\s{0,3}>\s?`)
	markdownListPattern         = regexp.MustCompile(`(?m)^\s*([-+*]|\d+\.)\s+`)
	markdownTokenPattern        = regexp.MustCompile("[*_`~]")
	whitespacePattern           = regexp.MustCompile(`\s+`)
)

func ExtractMoodRecordIDs(markdown string) ([]uint, error) {
	seen := make(map[uint]struct{})
	searchableMarkdown := stripCodeSections(markdown)

	for _, pattern := range []*regexp.Regexp{
		customMoodRecordLinkPattern,
		legacyMoodRecordLinkPattern,
	} {
		for _, match := range pattern.FindAllStringSubmatch(searchableMarkdown, -1) {
			if len(match) < 2 {
				continue
			}

			parsed, err := strconv.ParseUint(match[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse mood record link %q: %w", match[1], err)
			}
			seen[uint(parsed)] = struct{}{}
		}
	}

	ids := make([]uint, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	return ids, nil
}

func MarkdownPreview(markdown string) string {
	preview := strings.TrimSpace(markdown)
	if preview == "" {
		return ""
	}

	preview = customMoodRecordLinkPattern.ReplaceAllStringFunc(preview, moodRecordLinkPreviewText)
	preview = markdownLinkPattern.ReplaceAllString(preview, "$1")
	preview = markdownHeadingPattern.ReplaceAllString(preview, "")
	preview = markdownQuotePattern.ReplaceAllString(preview, "")
	preview = markdownListPattern.ReplaceAllString(preview, "")
	preview = markdownTokenPattern.ReplaceAllString(preview, "")
	preview = whitespacePattern.ReplaceAllString(preview, " ")
	preview = strings.TrimSpace(preview)

	if preview == "" {
		return ""
	}

	if utf8.RuneCountInString(preview) <= previewMaxLength {
		return preview
	}

	runes := []rune(preview)
	return strings.TrimSpace(string(runes[:previewMaxLength-1])) + "..."
}

func moodRecordLinkPreviewText(raw string) string {
	match := customMoodRecordLinkPattern.FindStringSubmatch(raw)
	if len(match) < 2 {
		return raw
	}

	if len(match) >= 3 && strings.TrimSpace(match[2]) != "" {
		return strings.TrimSpace(match[2])
	}

	return fmt.Sprintf("Mood record #%s", match[1])
}

func stripCodeSections(markdown string) string {
	withoutFencedCode := fencedCodeBlockPattern.ReplaceAllString(markdown, " ")
	return inlineCodePattern.ReplaceAllString(withoutFencedCode, " ")
}
