package journal_test

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/azaviyalov/null3/backend/internal/domain/journal"
)

func TestExtractMoodRecordIDs(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     []uint
	}{
		{
			name:     "custom and page links",
			markdown: `[[mood:42|Calm]] [[mood:7]] [mood](/mood-records/9) https://example.test/mood-records/12?view=full`,
			want:     []uint{7, 9, 12, 42},
		},
		{
			name:     "deduplicated and sorted",
			markdown: `[[mood:20]] /mood-records/3 [[mood:20|Again]] /mood-records/3#details`,
			want:     []uint{3, 20},
		},
		{
			name: "code spans and fences are ignored",
			markdown: strings.Join([]string{
				"[[mood:1]]",
				"`[[mood:2]]` and ``/mood-records/3``",
				"```markdown",
				"[[mood:4]] /mood-records/5",
				"```",
				"~~~",
				"[[mood:6]]",
				"~~~",
				"/mood-records/7",
			}, "\n"),
			want: []uint{1, 7},
		},
		{
			name:     "unclosed fence continues to end",
			markdown: "[[mood:1]]\n```\n[[mood:2]]\n/mood-records/3",
			want:     []uint{1},
		},
		{
			name:     "code span may cross lines",
			markdown: "`code starts\n[[mood:2]] /mood-records/3\ncode ends`\n[[mood:1]]",
			want:     []uint{1},
		},
		{
			name:     "ordinary text and malformed links",
			markdown: `mood:1 [[mood:nope]] /mood-records/nope`,
			want:     []uint{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := journal.ExtractMoodRecordIDs(tt.markdown)
			if err != nil {
				t.Fatalf("ExtractMoodRecordIDs() error = %v", err)
			}
			if !slices.Equal(got, tt.want) {
				t.Fatalf("ExtractMoodRecordIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractMoodRecordIDsReportsOverflow(t *testing.T) {
	_, err := journal.ExtractMoodRecordIDs(`[[mood:18446744073709551616]]`)

	if err == nil || !strings.Contains(err.Error(), "parse mood record link") {
		t.Fatalf("ExtractMoodRecordIDs() error = %v, want parsing context", err)
	}
	var numberError *strconv.NumError
	if !errors.As(err, &numberError) {
		t.Fatalf("ExtractMoodRecordIDs() error = %v, want strconv.NumError cause", err)
	}
}

func TestMarkdownPreview(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     string
	}{
		{
			name:     "markdown formatting and links",
			markdown: "# Heading\n\n> **Bold** [link](https://example.test)\n- item",
			want:     "Heading Bold link item",
		},
		{
			name:     "mood link labels",
			markdown: `[[mood:4| Peaceful ]] and [[mood:8]]`,
			want:     "Peaceful and Mood record #8",
		},
		{name: "blank", markdown: " \n\t ", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := journal.MarkdownPreview(tt.markdown); got != tt.want {
				t.Fatalf("MarkdownPreview() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMarkdownPreviewTruncatesByRunes(t *testing.T) {
	preview := journal.MarkdownPreview(strings.Repeat("я", 200))

	if !strings.HasSuffix(preview, "...") {
		t.Fatal("MarkdownPreview() does not end with an ellipsis")
	}
	if got := utf8.RuneCountInString(preview); got != 182 {
		t.Fatalf("MarkdownPreview() rune count = %d, want 182", got)
	}
}
