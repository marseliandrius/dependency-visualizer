package tests

import (
	"slices"
	"testing"

	"github.com/marseliandrius/dependency-visualizer/src/graph"
)

// TestOrder проверяет общий потомок, цепочку и пакет без зависимостей.
func TestOrder(t *testing.T) {
	source, err := graph.Load("repository-order.txt")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name string
		root string
		want []string
	}{
		{"shared", "A", []string{"D", "B", "E", "C", "A"}},
		{"chain", "B", []string{"D", "B"}},
		{"single", "F", []string{"F"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := source.Order(tc.root)
			if err != nil || !slices.Equal(got, tc.want) {
				t.Fatalf("получено %v, %v; ожидалось %v", got, err, tc.want)
			}
		})
	}
}

// TestOrderCycle проверяет отказ от строгого порядка при цикле.
func TestOrderCycle(t *testing.T) {
	source, err := graph.Load("repository.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := source.Order("A"); err == nil {
		t.Fatal("ожидалась ошибка циклической зависимости")
	}
}
