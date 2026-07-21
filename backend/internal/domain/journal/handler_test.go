package journal_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/core/server"
	"github.com/azaviyalov/null3/backend/internal/domain/account"
	"github.com/azaviyalov/null3/backend/internal/domain/journal"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/azaviyalov/null3/backend/internal/testutil"
	"github.com/labstack/echo/v4"
)

const journalTestJWTSecret = "journal-test-signing-secret"

func TestMoodRecordHTTPContracts(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newJournalTestEnvironment(t)
	owner := createJournalUser(t, environment, "owner")
	other := createJournalUser(t, environment, "other")
	e, tokenService := newJournalTestServer(t, environment)
	ownerCookie := journalUserCookie(t, tokenService, owner.ID)
	otherCookie := journalUserCookie(t, tokenService, other.ID)

	unauthorized := serveJournalJSON(t, e, http.MethodGet, "/api/journal/mood-records", nil)
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("unauthorized list status = %d, want %d", unauthorized.Code, http.StatusUnauthorized)
	}
	invalidRequests := []struct {
		name string
		body string
	}{
		{name: "malformed JSON", body: `{"feeling":`},
		{name: "missing feeling", body: `{}`},
	}
	for _, tt := range invalidRequests {
		t.Run(tt.name, func(t *testing.T) {
			response := serveJournalJSON(t, e, http.MethodPost, "/api/journal/mood-records", tt.body, ownerCookie)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
			}
		})
	}

	createResponse := serveJournalJSON(t, e, http.MethodPost, "/api/journal/mood-records", journal.MoodEditRecordRequest{
		Feeling: "calm",
		Emoji:   "🙂",
		Note:    "quiet morning",
	}, ownerCookie)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", createResponse.Code, http.StatusCreated)
	}
	if contentType := createResponse.Header().Get(echo.HeaderContentType); !strings.HasPrefix(contentType, echo.MIMEApplicationJSON) {
		t.Errorf("create Content-Type = %q, want JSON", contentType)
	}
	var created journal.MoodRecordResponse
	decodeJournalResponse(t, createResponse, &created)
	if created.ID == 0 || created.UserID != owner.ID || created.Feeling != "calm" {
		t.Error("create response does not contain the created owner record")
	}

	invalidID := serveJournalJSON(t, e, http.MethodGet, "/api/journal/mood-records/not-an-id", nil, ownerCookie)
	if invalidID.Code != http.StatusBadRequest {
		t.Fatalf("invalid ID status = %d, want %d", invalidID.Code, http.StatusBadRequest)
	}
	for _, request := range []struct {
		method string
		body   any
	}{
		{method: http.MethodGet},
		{method: http.MethodPut, body: journal.MoodEditRecordRequest{Feeling: "changed"}},
		{method: http.MethodDelete},
	} {
		response := serveJournalJSON(t, e, request.method, fmt.Sprintf("/api/journal/mood-records/%d", created.ID), request.body, otherCookie)
		if response.Code != http.StatusNotFound {
			t.Fatalf("foreign %s status = %d, want %d", request.method, response.Code, http.StatusNotFound)
		}
	}

	if _, err := environment.service.CreateMoodRecord(t.Context(), owner.ID, journal.MoodEditRecordRequest{Feeling: "second"}); err != nil {
		t.Fatalf("create second mood record: %v", err)
	}
	listResponse := serveJournalJSON(t, e, http.MethodGet, "/api/journal/mood-records?limit=1&offset=0", nil, ownerCookie)
	if listResponse.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", listResponse.Code, http.StatusOK)
	}
	var page core.Page[journal.MoodRecord]
	decodeJournalResponse(t, listResponse, &page)
	if page.TotalCount != 2 || len(page.Items) != 1 {
		t.Fatalf("list page length = %d total = %d, want 1 and 2", len(page.Items), page.TotalCount)
	}

	recordPath := fmt.Sprintf("/api/journal/mood-records/%d", created.ID)
	deleteResponse := serveJournalJSON(t, e, http.MethodDelete, recordPath, nil, ownerCookie)
	if deleteResponse.Code != http.StatusOK {
		t.Fatalf("delete status = %d, want %d", deleteResponse.Code, http.StatusOK)
	}
	getDeletedResponse := serveJournalJSON(t, e, http.MethodGet, recordPath, nil, ownerCookie)
	if getDeletedResponse.Code != http.StatusOK {
		t.Fatalf("get deleted status = %d, want %d", getDeletedResponse.Code, http.StatusOK)
	}
	deletedListResponse := serveJournalJSON(t, e, http.MethodGet, "/api/journal/mood-records?deleted=true", nil, ownerCookie)
	var deletedPage core.Page[journal.MoodRecord]
	decodeJournalResponse(t, deletedListResponse, &deletedPage)
	if deletedPage.TotalCount != 1 || len(deletedPage.Items) != 1 || deletedPage.Items[0].ID != created.ID {
		t.Fatal("deleted list does not contain the deleted mood record")
	}
	restoreResponse := serveJournalJSON(t, e, http.MethodPost, recordPath+"/restore", nil, ownerCookie)
	if restoreResponse.Code != http.StatusOK {
		t.Fatalf("restore status = %d, want %d", restoreResponse.Code, http.StatusOK)
	}
}

func TestDiaryEntryHTTPContracts(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newJournalTestEnvironment(t)
	owner := createJournalUser(t, environment, "owner")
	other := createJournalUser(t, environment, "other")
	e, tokenService := newJournalTestServer(t, environment)
	ownerCookie := journalUserCookie(t, tokenService, owner.ID)
	otherCookie := journalUserCookie(t, tokenService, other.ID)
	mood, err := environment.service.CreateMoodRecord(t.Context(), owner.ID, journal.MoodEditRecordRequest{Feeling: "calm"})
	if err != nil {
		t.Fatalf("create mood record: %v", err)
	}

	invalidRequests := []struct {
		name string
		body any
	}{
		{name: "malformed JSON", body: `{"markdown":`},
		{name: "missing occurred at", body: map[string]any{"markdown": "entry"}},
		{name: "future occurred at", body: journal.DiaryEditEntryRequest{
			Markdown:   "entry",
			OccurredAt: timePointer(time.Date(2100, time.January, 1, 0, 0, 0, 0, time.UTC)),
		}},
		{name: "foreign mood", body: journal.DiaryEditEntryRequest{
			Markdown:   "[[mood:999999]]",
			OccurredAt: timePointer(time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)),
		}},
	}
	for _, tt := range invalidRequests {
		t.Run(tt.name, func(t *testing.T) {
			response := serveJournalJSON(t, e, http.MethodPost, "/api/journal/diary-entries", tt.body, ownerCookie)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
			}
		})
	}

	occurredAt := time.Date(2026, time.January, 2, 9, 0, 0, 0, time.UTC)
	markdown := fmt.Sprintf("A **day** with [[mood:%d|calm]]", mood.ID)
	createResponse := serveJournalJSON(t, e, http.MethodPost, "/api/journal/diary-entries", journal.DiaryEditEntryRequest{
		Title:      "Day",
		Markdown:   markdown,
		OccurredAt: &occurredAt,
	}, ownerCookie)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", createResponse.Code, http.StatusCreated)
	}
	var created journal.DiaryEntryResponse
	decodeJournalResponse(t, createResponse, &created)
	if created.ID == 0 || created.UserID != owner.ID || created.Preview != "A day with calm" {
		t.Fatal("create response does not contain the expected diary projection")
	}
	if len(created.ReferencedMoodRecords) != 1 || created.ReferencedMoodRecords[0].ID != mood.ID {
		t.Fatal("create response does not contain the referenced mood record")
	}
	moodResponse := serveJournalJSON(t, e, http.MethodGet, fmt.Sprintf("/api/journal/mood-records/%d", mood.ID), nil, ownerCookie)
	if moodResponse.Code != http.StatusOK {
		t.Fatalf("mood backlink status = %d, want %d", moodResponse.Code, http.StatusOK)
	}
	var moodProjection journal.MoodRecordResponse
	decodeJournalResponse(t, moodResponse, &moodProjection)
	if len(moodProjection.DiaryEntryLinks) != 1 || moodProjection.DiaryEntryLinks[0].ID != created.ID || moodProjection.DiaryEntryLinks[0].Preview != "A day with calm" {
		t.Fatal("mood response does not contain the diary backlink projection")
	}

	entryPath := fmt.Sprintf("/api/journal/diary-entries/%d", created.ID)
	invalidID := serveJournalJSON(t, e, http.MethodGet, "/api/journal/diary-entries/not-an-id", nil, ownerCookie)
	if invalidID.Code != http.StatusBadRequest {
		t.Fatalf("invalid ID status = %d, want %d", invalidID.Code, http.StatusBadRequest)
	}
	for _, request := range []struct {
		method string
		body   any
	}{
		{method: http.MethodGet},
		{method: http.MethodPut, body: journal.DiaryEditEntryRequest{Markdown: "changed", OccurredAt: &occurredAt}},
		{method: http.MethodDelete},
	} {
		response := serveJournalJSON(t, e, request.method, entryPath, request.body, otherCookie)
		if response.Code != http.StatusNotFound {
			t.Fatalf("foreign %s status = %d, want %d", request.method, response.Code, http.StatusNotFound)
		}
	}

	deleteResponse := serveJournalJSON(t, e, http.MethodDelete, entryPath, nil, ownerCookie)
	if deleteResponse.Code != http.StatusOK {
		t.Fatalf("delete status = %d, want %d", deleteResponse.Code, http.StatusOK)
	}
	deletedListResponse := serveJournalJSON(t, e, http.MethodGet, "/api/journal/diary-entries?deleted=true", nil, ownerCookie)
	var deletedPage core.Page[journal.DiaryEntry]
	decodeJournalResponse(t, deletedListResponse, &deletedPage)
	if deletedPage.TotalCount != 1 || len(deletedPage.Items) != 1 || deletedPage.Items[0].ID != created.ID {
		t.Fatal("deleted list does not contain the deleted diary entry")
	}

	if _, err := environment.service.DeleteMoodRecord(t.Context(), owner.ID, mood.ID); err != nil {
		t.Fatalf("delete referenced mood record: %v", err)
	}
	rejectedRestore := serveJournalJSON(t, e, http.MethodPost, entryPath+"/restore", nil, ownerCookie)
	if rejectedRestore.Code != http.StatusBadRequest {
		t.Fatalf("restore with deleted mood status = %d, want %d", rejectedRestore.Code, http.StatusBadRequest)
	}
	if _, err := environment.service.RestoreMoodRecord(t.Context(), owner.ID, mood.ID); err != nil {
		t.Fatalf("restore referenced mood record: %v", err)
	}
	restoreResponse := serveJournalJSON(t, e, http.MethodPost, entryPath+"/restore", nil, ownerCookie)
	if restoreResponse.Code != http.StatusOK {
		t.Fatalf("restore status = %d, want %d", restoreResponse.Code, http.StatusOK)
	}
}

func newJournalTestServer(t *testing.T, environment *journalTestEnvironment) (*echo.Echo, *session.Service) {
	t.Helper()

	testutil.DiscardLogs(t)

	tokenService := session.NewService(session.NewRepository(environment.database), session.Config{
		JWTSecret:              journalTestJWTSecret,
		JWTExpiration:          time.Hour,
		RefreshTokenExpiration: 7 * 24 * time.Hour,
	})
	accountService := account.NewService(account.NewRepository(environment.database), tokenService, account.Config{
		PasswordResetTokenExpiration: time.Hour,
	})
	validateUser := func(ctx context.Context, userID uint) error {
		_, err := accountService.GetUserByID(ctx, userID)
		return err
	}
	e := server.NewEchoServer(server.Config{})
	journal.RegisterRoutes(e, journal.NewHandler(environment.service), session.UserJWTMiddleware(tokenService, validateUser))
	return e, tokenService
}

func journalUserCookie(t *testing.T, service *session.Service, userID uint) *http.Cookie {
	t.Helper()

	token, err := service.GenerateUserAccessToken(userID)
	if err != nil {
		t.Fatalf("generate user access token: %v", err)
	}
	return &http.Cookie{Name: session.UserCookieName, Value: token}
}

func serveJournalJSON(t *testing.T, e *echo.Echo, method, path string, body any, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	t.Helper()
	return testutil.JSONRequest(t, e, method, path, body, cookies...)
}

func decodeJournalResponse(t *testing.T, response *httptest.ResponseRecorder, target any) {
	t.Helper()
	testutil.DecodeJSON(t, response, target)
}
