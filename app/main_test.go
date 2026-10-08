package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fairhive-labs/go-pwdgen/pkg/generator"
)

type response struct {
	Length   int
	Password string
}

func TestGenerate(t *testing.T) {
	tt := []struct {
		name   string
		length int
		url    string
	}{
		{"normal", 16, "/?mime=json"},
		{"16", 16, fmt.Sprintf("/?l=%d&mime=json", 16)},
		{"100", 100, fmt.Sprintf("/?l=%d&mime=json", 100)},
		{"32", 32, fmt.Sprintf("/?l=%d&mime=json", 32)},
		{"min length 8->10", 10, fmt.Sprintf("/?l=%d&mime=json", 8)},
		{"max length", generator.MaxLength, fmt.Sprintf("/?l=%d&mime=json", generator.MaxLength)},
		{"max length + 1", generator.MinLength, fmt.Sprintf("/?l=%d&mime=json", generator.MaxLength+1)},
		{"incorrect", 16, fmt.Sprintf("/?l=%s&mime=json", "foo")},
	}

	for _, tc := range tt {
		t.Run("json_"+tc.name, func(t *testing.T) {
			router := setupRouter()
			w := httptest.NewRecorder()
			l := tc.length
			req, _ := http.NewRequest("GET", tc.url, nil)
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("%d should be %d", w.Code, http.StatusOK)
				t.FailNow()
			}

			var r response

			json.NewDecoder(w.Body).Decode(&r)

			if r.Length != l {
				t.Errorf("incorrect length, got %d, want %d", r.Length, l)
				t.FailNow()
			}

			if len(r.Password) != l {
				t.Errorf("incorrect length of %q, got %d, want %d", r.Password, len(r.Password), l)
				t.FailNow()
			}
		})
	}

	tt = []struct {
		name   string
		length int
		url    string
	}{
		{"no mime l 20", 20, "/?l=20"},
		{"mime l 32", 32, "/?l=32&mime=html"},
		{"incorrect mime type", 16, "/?mime=foo"},
	}
	for _, tc := range tt {
		t.Run("html_"+tc.name, func(t *testing.T) {
			router := setupRouter()
			w := httptest.NewRecorder()
			l := tc.length
			req, _ := http.NewRequest("GET", tc.url, nil)
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("%d should be %d", w.Code, http.StatusOK)
				t.FailNow()
			}

			if w.Body.Len() == 0 || w.Body == nil {
				t.Errorf("response body cannot be empty or nil")
				t.FailNow()
			}

			m := map[string]bool{
				"PWDGEN":                             false,
				"SECURE FORGE ONLINE":                false,
				fmt.Sprintf("data-length=\"%d\"", l): false,
			}

			for l, err := w.Body.ReadString('\n'); err == nil; {
				for k := range m {
					if strings.Contains(l, k) {
						m[k] = true
					}
				}
				l, err = w.Body.ReadString('\n')
			}

			for k, v := range m {
				if !v {
					t.Errorf("response body should contain %q", k)
					t.FailNow()
				}
			}
		})
	}
}

// Controls the JSON payload: content type, exact keys, and a password that
// changes on every request.
func TestGenerateJSONContract(t *testing.T) {
	router := setupRouter()
	seen := map[string]bool{}

	for range 5 {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/?l=24&mime=json", nil)
		router.ServeHTTP(w, req)

		if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
			t.Fatalf("incorrect content type, got %q, want application/json", ct)
		}

		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("response is not valid JSON: %v", err)
		}
		if len(body) != 2 {
			t.Errorf("response should only hold length and password, got %v", body)
		}
		pwd, ok := body["password"].(string)
		if !ok || len(pwd) != 24 {
			t.Fatalf("incorrect password %v, want a 24-char string", body["password"])
		}
		if l, ok := body["length"].(float64); !ok || l != 24 {
			t.Errorf("incorrect length %v, want 24", body["length"])
		}
		if seen[pwd] {
			t.Errorf("password %q was served twice", pwd)
		}
		seen[pwd] = true
	}
}

// Controls the HTML page embeds the generated password.
func TestGenerateHTMLContract(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?l=20", nil)
	setupRouter().ServeHTTP(w, req)

	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("incorrect content type, got %q, want text/html", ct)
	}
	if !strings.Contains(w.Body.String(), `data-length="20"`) {
		t.Errorf("page should expose the password length")
	}
}

// Controls lengths that are not usable fall back instead of failing.
func TestGenerateFallbackLength(t *testing.T) {
	tt := []struct {
		name   string
		query  string
		length int
	}{
		{"negative", "l=-5", generator.MinLength},
		{"zero", "l=0", generator.MinLength},
		{"empty", "l=", length},
		{"float", "l=12.5", length},
		{"overflow", "l=99999999999999999999", length},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/?mime=json&"+tc.query, nil)
			setupRouter().ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("%d should be %d", w.Code, http.StatusOK)
			}
			var r response
			if err := json.NewDecoder(w.Body).Decode(&r); err != nil {
				t.Fatalf("response is not valid JSON: %v", err)
			}
			if r.Length != tc.length || len(r.Password) != tc.length {
				t.Errorf("incorrect length, got %d (%q), want %d", r.Length, r.Password, tc.length)
			}
		})
	}
}

// Controls only GET / is served.
func TestRoutes(t *testing.T) {
	tt := []struct {
		method string
		path   string
		code   int
	}{
		{"GET", "/", http.StatusOK},
		{"GET", "/unknown", http.StatusNotFound},
		{"POST", "/", http.StatusNotFound},
	}
	for _, tc := range tt {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			setupRouter().ServeHTTP(w, req)

			if w.Code != tc.code {
				t.Errorf("got %d, want %d", w.Code, tc.code)
			}
		})
	}
}
