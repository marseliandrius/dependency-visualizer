package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/marseliandrius/dependency-visualizer/src/config"
)

// main читает параметры запуска и выводит конфигурацию первого этапа.
func main() {
	cfg, err := config.Parse(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}

	fmt.Printf("package=%q\n", cfg.Package)
	fmt.Printf("repository=%q\n", cfg.Repository)
	fmt.Printf("test=%t\n", cfg.TestMode)
	fmt.Printf("version=%q\n", cfg.Version)
	fmt.Printf("output=%q\n", cfg.Output)
	fmt.Printf("max-depth=%d\n", cfg.MaxDepth)
	fmt.Printf("filter=%q\n", cfg.Filter)
}
