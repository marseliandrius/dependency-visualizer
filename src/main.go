package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/marseliandrius/dependency-visualizer/src/apt"
	"github.com/marseliandrius/dependency-visualizer/src/config"
)

const noDependencies = 0

// main читает параметры и выводит прямые зависимости заданной версии пакета.
func main() {
	cfg, err := config.Parse(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}

	if cfg.TestMode {
		fmt.Fprintln(os.Stderr, "Ошибка: режим тестового графа будет реализован на этапе 3")
		os.Exit(1)
	}
	dependencies, err := apt.Fetch(cfg.Repository, cfg.Package, cfg.Version)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
	if len(dependencies) == noDependencies {
		fmt.Println("Прямых зависимостей нет.")
		return
	}
	for _, dependency := range dependencies {
		fmt.Println(dependency)
	}
}
