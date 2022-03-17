package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type response struct {
	Length   int
	Password string
}

func TestJSONRoute(t *testing.T) {
	tt := []struct {
		name   string
		length int
		url    string
	}{
		{"normal", 16, "/json"},
		{"16", 16, fmt.Sprintf("/json?l=%d", 16)},
		{"100", 100, fmt.Sprintf("/json?l=%d", 100)},
		{"32", 32, fmt.Sprintf("/json?l=%d", 32)},
		{"too short 8->10", 10, fmt.Sprintf("/json?l=%d", 8)},
		{"incorrect", 16, fmt.Sprintf("/json?l=%s", "foo")},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupRouter("../templates/*")
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
				t.Errorf("length should be %d", l)
				t.FailNow()
			}

			if len(r.Password) != l {
				t.Errorf("%s length is %d and should be %d", r.Password, len(r.Password), l)
				t.FailNow()
			}
		})
	}
}

func TestHTMLRoute(t *testing.T) {
	router := setupRouter("../templates/*")
	w := httptest.NewRecorder()
	l := 20
	req, _ := http.NewRequest("GET", fmt.Sprintf("/?l=%d", l), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("%d should be %d", w.Code, http.StatusOK)
		t.FailNow()
	}

	if w.Body.Len() == 0 || w.Body == nil {
		t.Errorf("response body cannot be empty or nil")
		t.FailNow()
	}

	t.Log(w.Body)

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
}
