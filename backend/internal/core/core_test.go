package core_test

import (
	"encoding/json"
	"testing"

	"github.com/azaviyalov/null3/backend/internal/core"
)

func TestDeletedFilterModeString(t *testing.T) {
	tests := []struct {
		name string
		mode core.DeletedFilterMode
		want string
	}{
		{name: "non-deleted", mode: core.DeletedModeNonDeleted, want: "NonDeleted"},
		{name: "deleted only", mode: core.DeletedModeDeletedOnly, want: "DeletedOnly"},
		{name: "all", mode: core.DeletedModeAll, want: "All"},
		{name: "unknown negative value", mode: core.DeletedFilterMode(-1), want: "Unknown"},
		{name: "unknown positive value", mode: core.DeletedFilterMode(99), want: "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.want {
				t.Fatalf("DeletedFilterMode(%d).String() = %q, want %q", tt.mode, got, tt.want)
			}
		})
	}
}

func TestPageJSONContract(t *testing.T) {
	page := core.Page[string]{
		Items:      []string{"first", "second"},
		TotalCount: 7,
	}

	data, err := json.Marshal(page)
	if err != nil {
		t.Fatalf("marshal page: %v", err)
	}

	const want = `{"items":["first","second"],"total_count":7}`
	if string(data) != want {
		t.Fatalf("page JSON = %s, want %s", data, want)
	}
}
