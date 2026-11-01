package tests

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/marseliandrius/dependency-visualizer/src/apt"
)

// TestFind проверяет точную версию, отсутствие зависимостей и неизвестную версию.
func TestFind(t *testing.T) {
	data, err := os.ReadFile("Packages")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name, version string
		want          []string
		wantError     bool
	}{
		{"demo", "2.0", []string{"base (>= 1.0)", "alpha:any (>= 2.0)", "beta | gamma"}, false},
		{"nodeps", "1.0", nil, false},
		{"demo", "9.0", nil, true},
	}
	for _, tc := range cases {
		t.Run(tc.name+"/"+tc.version, func(t *testing.T) {
			got, err := apt.Find(bytes.NewReader(data), tc.name, tc.version)
			if (err != nil) != tc.wantError || !slices.Equal(got, tc.want) {
				t.Fatalf("получено %v, %v; ожидалось %v, ошибка=%t", got, err, tc.want, tc.wantError)
			}
		})
	}
}

// TestFetch проверяет HTTP-загрузку обычного индекса, gzip и ошибку 404.
func TestFetch(t *testing.T) {
	data := []byte("Package: demo\nVersion: 1.0\nDepends: alpha\n")
	var compressed bytes.Buffer
	archive := gzip.NewWriter(&compressed)
	if _, err := archive.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/missing" {
			http.NotFound(w, r)
			return
		}
		body := data
		if strings.HasSuffix(r.URL.Path, ".gz") {
			body = compressed.Bytes()
		}
		if _, err := w.Write(body); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	for _, path := range []string{"/Packages", "/Packages.gz"} {
		got, err := apt.Fetch(server.URL+path, "demo", "1.0")
		if err != nil || !slices.Equal(got, []string{"alpha"}) {
			t.Fatalf("%s: получено %v, %v", path, got, err)
		}
	}
	if _, err := apt.Fetch(server.URL+"/missing", "demo", "1.0"); err == nil {
		t.Fatal("ожидалась ошибка HTTP 404")
	}
}
