package diary

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
	customMoodEntryLinkPattern = regexp.MustCompile(`\[\[mood:(\d+)(?:\|([^\]]+))?\]\]`)
	legacyMoodEntryLinkPattern = regexp.MustCompile(`(?:https?://[^\s)]+)?/mood/entries/(\d+)(?:[?#][^\s)]*)?`)
	fencedCodeBlockPattern     = regexp.MustCompile("(?s)```.*?```|~~~.*?~~~")
	inlineCodePattern          = regexp.MustCompile("`[^`\n]*`")
	markdownLinkPattern        = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	markdownHeadingPattern     = regexp.MustCompile(`(?m)^\s{0,3}#{1,6}\s*`)
	markdownQuotePattern       = regexp.MustCompile(`(?m)^\s{0,3}>\s?`)
	markdownListPattern        = regexp.MustCompile(`(?m)^\s*([-+*]|\d+\.)\s+`)
	markdownTokenPattern       = regexp.MustCompile("[*_`~]")
	whitespacePattern          = regexp.MustCompile(`\s+`)
)

func ExtractMoodEntryIDs(markdown string) ([]uint, error) {
	seen := make(map[uint]struct{})
	searchableMarkdown := stripCodeSections(markdown)

	for _, pattern := range []*regexp.Regexp{
		customMoodEntryLinkPattern,
		legacyMoodEntryLinkPattern,
	} {
		for _, match := range pattern.FindAllStringSubmatch(searchableMarkdown, -1) {
			if len(match) < 2 {
				continue
			}

			parsed, err := strconv.ParseUint(match[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse mood entry link %q: %w", match[1], err)
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

	preview = customMoodEntryLinkPattern.ReplaceAllStringFunc(preview, moodEntryLinkPreviewText)
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

func moodEntryLinkPreviewText(raw string) string {
	match := customMoodEntryLinkPattern.FindStringSubmatch(raw)
	if len(match) < 2 {
		return raw
	}

	if len(match) >= 3 && strings.TrimSpace(match[2]) != "" {
		return strings.TrimSpace(match[2])
	}

	return fmt.Sprintf("Mood Entry #%s", match[1])
}

func stripCodeSections(markdown string) string {
	withoutFencedCode := fencedCodeBlockPattern.ReplaceAllString(markdown, " ")
	return inlineCodePattern.ReplaceAllString(withoutFencedCode, " ")
}
