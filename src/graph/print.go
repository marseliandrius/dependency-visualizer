package graph

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

// Print выводит списки смежности, сортируя имена вершин для устойчивого результата.
func (g Graph) Print(writer io.Writer) error {
	names := make([]string, 0, len(g))
	for name := range g {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		line := name + ":"
		if dependencies := strings.Join(g[name], " "); dependencies != "" {
			line += " " + dependencies
		}
		if _, err := fmt.Fprintln(writer, line); err != nil {
			return err
		}
	}
	return nil
}
