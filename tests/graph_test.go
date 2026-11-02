package tests

import (
	"slices"
	"strings"
	"testing"

	"github.com/marseliandrius/dependency-visualizer/src/apt"
	"github.com/marseliandrius/dependency-visualizer/src/graph"
)

// TestBuild проверяет цикл, меньшую глубину повторного пути, корень и фильтрацию.
func TestBuild(t *testing.T) {
	source, err := graph.Load("repository.txt")
	if err != nil {
		t.Fatal(err)
	}
	full := graph.Graph{"A": {"B", "D"}, "B": {"C"}, "C": {"D"}, "D": {"E"}, "E": {"F"}, "F": {"A"}}
	limited := graph.Graph{"A": {"B", "D"}, "B": {"C"}, "C": {"D"}, "D": {"E"}, "E": {"F"}, "F": nil}
	cases := []struct {
		name     string
		maxDepth int
		filter   string
		want     graph.Graph
	}{
		{"cycle", 6, "", full},
		{"depth", 3, "", limited},
		{"root", 0, "", graph.Graph{"A": nil}},
		{"filter", 6, "D", graph.Graph{"A": {"B"}, "B": {"C"}, "C": nil}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := graph.Build(source, "A", tc.maxDepth, tc.filter)
			if err != nil || !sameGraph(got, tc.want) {
				t.Fatalf("получено %v, %v; ожидалось %v", got, err, tc.want)
			}
		})
	}
}

// TestReadInvalid проверяет запрет строчных букв в тестовом репозитории.
func TestReadInvalid(t *testing.T) {
	if _, err := graph.Read(strings.NewReader("a: B\nB:\n")); err == nil {
		t.Fatal("ожидалась ошибка имени пакета")
	}
}

// TestReadGraph проверяет версию корня, альтернативы и поставщика виртуального пакета.
func TestReadGraph(t *testing.T) {
	data := "Package: demo\nVersion: 1.0\nDepends: old\n\n" +
		"Package: demo\nVersion: 2.0\nDepends: alpha:any(>= 2.0), missing | abstract\n\n" +
		"Package: alpha\nVersion: 2.0\n\nPackage: nodeps\nVersion: 1.0\nProvides: abstract\n"
	got, err := apt.ReadGraph(strings.NewReader(data), "demo", "2.0")
	want := graph.Graph{"demo": {"alpha", "nodeps"}, "alpha": nil, "nodeps": nil}
	if err != nil || !sameGraph(got, want) {
		t.Fatalf("получено %v, %v; ожидалось %v", got, err, want)
	}
}

// sameGraph сравнивает наборы вершин и списки зависимостей без отражения типов.
func sameGraph(got, want graph.Graph) bool {
	if len(got) != len(want) {
		return false
	}
	for name, dependencies := range want {
		actual, exists := got[name]
		if !exists || !slices.Equal(actual, dependencies) {
			return false
		}
	}
	return true
}
