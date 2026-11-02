package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/marseliandrius/dependency-visualizer/src/config"
	"github.com/marseliandrius/dependency-visualizer/src/graph"
)

const emptyGraph = 0

// main запускает построение графа и сообщает об ошибке выполнения.
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}

// run получает настройки, загружает репозиторий и строит граф зависимостей.
func run() error {
	cfg, err := config.Parse(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return nil
	}
	if err != nil {
		return err
	}
	source, err := loadRepository(cfg)
	if err != nil {
		return err
	}
	result, err := graph.Build(source, cfg.Package, cfg.MaxDepth, cfg.Filter)
	if err != nil {
		return err
	}
	if len(result) == emptyGraph {
		fmt.Println("Граф пуст: исходный пакет исключён фильтром.")
		return nil
	}
	return result.Print(os.Stdout)
}
