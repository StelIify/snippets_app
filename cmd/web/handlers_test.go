package main

import (
	"net/http"
	"testing"
	"net/http/httptest"
)

func TestHomeHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(home)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    expected := "Hello from SnippetBox"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}

func TestSnippetCreateHandler(t *testing.T) {
	// Test POST request with /snippet/create URL
	req, err := http.NewRequest("POST", "/snippet/create", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Test GET request with /snippet/create URL
	reqGet, err := http.NewRequest("GET", "/snippet/create", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Test the handler for POST request
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(snippetCreate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Hello from snippet create"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Test the handler for non-POST request
	rrGet := httptest.NewRecorder()
	handler.ServeHTTP(rrGet, reqGet)

	if status := rrGet.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}

	if allow := rrGet.Header().Get("Allow"); allow != "POST" {
		t.Errorf("handler returned wrong Allow header: got %v want %v", allow, "POST")
	}
}