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
