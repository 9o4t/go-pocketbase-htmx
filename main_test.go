package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase"
)

func TestHealthzAndNotesRoundTrip(t *testing.T) {
	app := pocketbase.NewWithConfig(pocketbase.Config{DefaultDataDir: t.TempDir()})
	if err := app.Bootstrap(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = app.ResetBootstrapState() })
	if err := ensureNotes(app); err != nil {
		t.Fatal(err)
	}

	h := handler(app)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "ok" {
		t.Fatalf("healthz %d %q", rec.Code, rec.Body.String())
	}

	home := httptest.NewRecorder()
	h.ServeHTTP(home, httptest.NewRequest(http.MethodGet, "/", nil))
	if home.Code != http.StatusOK || !strings.Contains(home.Body.String(), "hx-post=\"/notes\"") {
		t.Fatalf("home %d %s", home.Code, home.Body.String())
	}

	form := url.Values{"title": {"hello"}, "body": {"from htmx"}}
	post := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(form.Encode()))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	out := httptest.NewRecorder()
	h.ServeHTTP(out, post)
	if out.Code != http.StatusOK || !strings.Contains(out.Body.String(), "hello") {
		t.Fatalf("post %d %s", out.Code, out.Body.String())
	}
}
