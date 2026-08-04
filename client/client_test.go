package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akbariandev/quickwit-gosdk/types"
)

func TestNewWithOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	c := New(ts.URL, WithTimeout(30*time.Second))
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	if c.HTTP == nil {
		t.Fatal("expected non-nil HTTP client")
	}
}

func TestNewWithTransport(t *testing.T) {
	var gotRequest bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRequest = true
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	c := New(ts.URL, WithTransport(ts.Client().Transport))
	if _, err := c.HTTP.R().Get("/"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !gotRequest {
		t.Error("expected request to reach server")
	}
}

func TestErrorTypeFormatting(t *testing.T) {
	e := &types.Error{StatusCode: 404, Message: "index not found"}
	if e.Error() != "quickwit api error (status 404): index not found" {
		t.Errorf("unexpected error message: %s", e.Error())
	}
}

func TestErrorHookSurfacesHTTPError(t *testing.T) {
	// Non-2xx responses must be converted into types.Error, not silently ignored.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"index already exists"}`))
	}))
	defer ts.Close()

	c := New(ts.URL)
	_, err := c.HTTP.R().Get("/")
	if err == nil {
		t.Fatal("expected error for 400 response, got nil")
	}
	te, ok := err.(*types.Error)
	if !ok {
		t.Fatalf("expected *types.Error, got %T: %v", err, err)
	}
	if te.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", te.StatusCode)
	}
	if te.Message != "index already exists" {
		t.Errorf("expected message 'index already exists', got %q", te.Message)
	}
}

func TestErrorHookNoBody(t *testing.T) {
	// Non-2xx with empty/unparseable body should still return an error with the status.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := New(ts.URL)
	_, err := c.HTTP.R().Get("/")
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
	te, ok := err.(*types.Error)
	if !ok {
		t.Fatalf("expected *types.Error, got %T: %v", err, err)
	}
	if te.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", te.StatusCode)
	}
}
