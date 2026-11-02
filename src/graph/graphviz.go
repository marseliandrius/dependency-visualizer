package graph

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

// DOT формирует описание ориентированного графа на языке Graphviz.
func (g Graph) DOT() string {
	names := make([]string, 0, len(g))
	for name := range g {
		names = append(names, name)
	}
	slices.Sort(names)
	var text strings.Builder
	fmt.Fprintln(&text, "digraph dependencies {")
	for _, name := range names {
		fmt.Fprintf(&text, "  %q;\n", name)
		for _, dependency := range g[name] {
			fmt.Fprintf(&text, "  %q -> %q;\n", name, dependency)
		}
	}
	fmt.Fprintln(&text, "}")
	return text.String()
}

// SavePNG передаёт описание графа программе dot и сохраняет PNG по заданному пути.
func (g Graph) SavePNG(path string) error {
	cmd := exec.Command("dot", "-Tpng", "-o", path)
	cmd.Stdin = strings.NewReader(g.DOT())
	message, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("создание PNG через Graphviz (dot): %w: %s", err, strings.TrimSpace(string(message)))
	}
	return nil
}
