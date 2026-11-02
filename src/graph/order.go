package graph

import "fmt"

const (
	unvisited = iota
	visiting
	finished
)

// ordering хранит состояния вершин и результат обхода для порядка загрузки.
type ordering struct {
	source Graph
	states map[string]int
	result []string
}

// Order возвращает зависимости перед зависящими от них пакетами, включая корень.
func (g Graph) Order(root string) ([]string, error) {
	o := ordering{source: g, states: make(map[string]int)}
	if err := o.visit(root); err != nil {
		return nil, err
	}
	return o.result, nil
}

// visit добавляет вершину после её зависимостей и обнаруживает обратные рёбра.
func (o *ordering) visit(name string) error {
	switch o.states[name] {
	case visiting:
		return fmt.Errorf("цикл через пакет %s: строгий порядок загрузки невозможен", name)
	case finished:
		return nil
	}
	dependencies, exists := o.source[name]
	if !exists {
		return fmt.Errorf("пакет %s отсутствует в графе", name)
	}
	o.states[name] = visiting
	for _, dependency := range dependencies {
		if err := o.visit(dependency); err != nil {
			return err
		}
	}
	o.states[name] = finished
	o.result = append(o.result, name)
	return nil
}
