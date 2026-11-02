package tests

import (
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/marseliandrius/dependency-visualizer/src/graph"
)

// TestDOT проверяет отдельную вершину, цикл и имена пакетов с дефисами.
func TestDOT(t *testing.T) {
	cases := []struct {
		name   string
		source graph.Graph
		want   string
	}{
		{"single", graph.Graph{"A": nil}, "digraph dependencies {\n  \"A\";\n}\n"},
		{"cycle", graph.Graph{"demo": {"lib-a"}, "lib-a": {"demo"}},
			"digraph dependencies {\n  \"demo\";\n  \"demo\" -> \"lib-a\";\n" +
				"  \"lib-a\";\n  \"lib-a\" -> \"demo\";\n}\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.source.DOT(); got != tc.want {
				t.Fatalf("получено %q; ожидалось %q", got, tc.want)
			}
		})
	}
}

// TestSavePNG проверяет создание настоящего PNG с циклическими зависимостями.
func TestSavePNG(t *testing.T) {
	if _, err := exec.LookPath("dot"); err != nil {
		t.Skip("для проверки PNG требуется Graphviz (dot) в PATH")
	}
	path := filepath.Join(t.TempDir(), "graph.png")
	source := graph.Graph{"A": {"B"}, "B": {"A"}}
	if err := source.SavePNG(path); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if _, err := png.Decode(file); err != nil {
		t.Fatal(err)
	}
}
