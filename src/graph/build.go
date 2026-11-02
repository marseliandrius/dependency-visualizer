// Package graph загружает, строит и выводит граф зависимостей.
package graph

import (
	"fmt"
	"strings"
)

const rootDepth = 0

// Graph хранит списки смежности: имя пакета и имена его прямых зависимостей.
type Graph map[string][]string

// walker хранит состояние рекурсивного DFS в пределах одного построения графа.
type walker struct {
	source   Graph
	result   Graph
	depths   map[string]int
	maxDepth int
	filter   string
}

// Build строит достижимую часть графа, учитывая глубину, фильтр и циклы.
func Build(source Graph, root string, maxDepth int, filter string) (Graph, error) {
	if maxDepth < rootDepth {
		return nil, fmt.Errorf("глубина не может быть отрицательной")
	}
	w := walker{
		source: source, result: make(Graph), depths: make(map[string]int),
		maxDepth: maxDepth, filter: filter,
	}
	if err := w.visit(root, rootDepth); err != nil {
		return nil, err
	}
	return w.result, nil
}

// visit рекурсивно раскрывает вершину, повторяя обработку лишь при меньшей глубине.
func (w *walker) visit(name string, depth int) error {
	if w.excluded(name) {
		return nil
	}
	previous, visited := w.depths[name]
	if visited && previous <= depth {
		return nil
	}
	dependencies, exists := w.source[name]
	if !exists {
		return fmt.Errorf("пакет %s отсутствует в репозитории", name)
	}
	w.depths[name] = depth
	w.result[name] = nil
	if depth == w.maxDepth {
		return nil
	}
	for _, dependency := range dependencies {
		if w.excluded(dependency) {
			continue
		}
		w.result[name] = append(w.result[name], dependency)
		if err := w.visit(dependency, depth+1); err != nil {
			return err
		}
	}
	return nil
}

// excluded проверяет совпадение имени с непустым фильтром.
func (w *walker) excluded(name string) bool {
	return w.filter != "" && strings.Contains(name, w.filter)
}
