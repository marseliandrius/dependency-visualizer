package tests

import (
	"strings"
	"testing"

	"github.com/marseliandrius/dependency-visualizer/src/config"
)

// TestParseValid проверяет URL, значения по умолчанию и все опции тестового режима.
func TestParseValid(t *testing.T) {
	args := []string{"-package=bash", "-repository=https://example.org/Packages", "-version=1.0"}
	cfg, err := config.Parse(args)
	want := config.Config{
		Package: "bash", Repository: "https://example.org/Packages", Version: "1.0",
		Output: "graph.png", MaxDepth: 3,
	}
	if err != nil || cfg != want {
		t.Fatalf("получено %+v, %v; ожидалось %+v без ошибки", cfg, err, want)
	}
	args = []string{
		"-package=A", "-repository=repository.txt", "-test=true", "-version=1.0",
		"-output=test.png", "-max-depth=0", "-filter=B",
	}
	cfg, err = config.Parse(args)
	want = config.Config{
		Package: "A", Repository: "repository.txt", TestMode: true, Version: "1.0",
		Output: "test.png", MaxDepth: 0, Filter: "B",
	}
	if err != nil || cfg != want {
		t.Fatalf("получено %+v, %v; ожидалось %+v без ошибки", cfg, err, want)
	}
}

// TestParseErrors проверяет по одной ошибке для каждого из семи параметров.
func TestParseErrors(t *testing.T) {
	cases := []struct {
		name string
		arg  string
	}{
		{"package", "-package="},
		{"repository", "-repository=bad-url"},
		{"test", "-test=maybe"},
		{"version", "-version="},
		{"output", "-output=graph.txt"},
		{"max-depth", "-max-depth=-1"},
		{"filter", "-filter=bad value"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			args := []string{"-package=bash", "-repository=https://example.org/Packages", "-version=1.0"}
			args = append(args, tc.arg)
			_, err := config.Parse(args)
			if err == nil || !strings.Contains(err.Error(), tc.name) {
				t.Fatalf("ожидалась ошибка параметра %s, получено: %v", tc.name, err)
			}
		})
	}
}
