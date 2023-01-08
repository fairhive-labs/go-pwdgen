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
				"Length":   false,
				"Password": false,
				fmt.Sprintf("<td><code>%d</code></td>", l): false,
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
