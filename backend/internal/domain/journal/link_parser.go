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
	moodRecordPageLinkPattern   = regexp.MustCompile(`(?:https?://[^\s)]+)?/mood-records/(\d+)(?:[?#][^\s)]*)?`)
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
		moodRecordPageLinkPattern,
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
	var result strings.Builder
	inFence := false
	var fenceCharacter byte
	fenceLength := 0

	for _, rawLine := range strings.SplitAfter(markdown, "\n") {
		line, hasNewline := strings.CutSuffix(rawLine, "\n")
		if inFence {
			if isClosingFence(line, fenceCharacter, fenceLength) {
				inFence = false
			}
			if hasNewline {
				result.WriteByte('\n')
			} else {
				result.WriteByte(' ')
			}
			continue
		}

		if character, length, ok := openingFence(line); ok {
			inFence = true
			fenceCharacter = character
			fenceLength = length
			if hasNewline {
				result.WriteByte('\n')
			} else {
				result.WriteByte(' ')
			}
			continue
		}

		result.WriteString(line)
		if hasNewline {
			result.WriteByte('\n')
		}
	}

	return stripInlineCode(result.String())
}

func openingFence(line string) (byte, int, bool) {
	line, ok := trimFenceIndent(line)
	if !ok || len(line) < 3 || (line[0] != '`' && line[0] != '~') {
		return 0, 0, false
	}

	character := line[0]
	length := countRun(line, character)
	if length < 3 || character == '`' && strings.ContainsRune(line[length:], '`') {
		return 0, 0, false
	}
	return character, length, true
}

func isClosingFence(line string, character byte, minimumLength int) bool {
	line, ok := trimFenceIndent(line)
	if !ok || len(line) == 0 || line[0] != character {
		return false
	}
	length := countRun(line, character)
	return length >= minimumLength && strings.TrimSpace(line[length:]) == ""
}

func trimFenceIndent(line string) (string, bool) {
	indent := 0
	for indent < len(line) && line[indent] == ' ' {
		indent++
	}
	if indent > 3 {
		return "", false
	}
	return line[indent:], true
}

func stripInlineCode(markdown string) string {
	var result strings.Builder
	for position := 0; position < len(markdown); {
		start := strings.IndexByte(markdown[position:], '`')
		if start < 0 {
			result.WriteString(markdown[position:])
			break
		}
		start += position
		result.WriteString(markdown[position:start])

		length := countRun(markdown[start:], '`')
		end := findBacktickRun(markdown, start+length, length)
		if end < 0 {
			result.WriteString(markdown[start : start+length])
			position = start + length
			continue
		}

		result.WriteByte(' ')
		position = end + length
	}
	return result.String()
}

func findBacktickRun(line string, position, length int) int {
	for position < len(line) {
		index := strings.IndexByte(line[position:], '`')
		if index < 0 {
			return -1
		}
		index += position
		runLength := countRun(line[index:], '`')
		if runLength == length {
			return index
		}
		position = index + runLength
	}
	return -1
}

func countRun(value string, character byte) int {
	length := 0
	for length < len(value) && value[length] == character {
		length++
	}
	return length
}
