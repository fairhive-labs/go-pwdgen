package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type response struct {
	Length   int
	Password string
}

func TestMainRoute(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("%d should be %d", w.Code, http.StatusOK)
		t.FailNow()
	}

	length := 16
	var r response

	json.NewDecoder(w.Body).Decode(&r)

	if r.Length != length {
		t.Errorf("length should be %d", length)
		t.FailNow()
	}

	if len(r.Password) != length {
		t.Errorf("%s length should be %d", r.Password, length)
		t.FailNow()
	}

}

func TestMainRouteWith32Length(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	length := 32
	req, _ := http.NewRequest("GET", fmt.Sprintf("/?l=%d", length), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("%d should be %d", w.Code, http.StatusOK)
		t.FailNow()
	}

	var r response

	json.NewDecoder(w.Body).Decode(&r)

	if r.Length != length {
		t.Errorf("length should be %d", length)
		t.FailNow()
	}

	if len(r.Password) != length {
		t.Errorf("%s length should be %d", r.Password, length)
		t.FailNow()
	}
}

func TestMainRouteWithIncompatibleLength(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	length := 16
	req, _ := http.NewRequest("GET", fmt.Sprintf("/?l=%s", "mistake"), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("%d should be %d", w.Code, http.StatusOK)
		t.FailNow()
	}

	var r response

	json.NewDecoder(w.Body).Decode(&r)

	if r.Length != length {
		t.Errorf("length should be %d", length)
		t.FailNow()
	}

	if len(r.Password) != length {
		t.Errorf("%s length should be %d", r.Password, length)
		t.FailNow()
	}
}
