// Package config разбирает и проверяет параметры командной строки.
package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultMaxDepth = 3
	minDepth        = 0
	noArguments     = 0
)

// Config хранит семь настраиваемых параметров приложения.
type Config struct {
	Package    string
	Repository string
	TestMode   bool
	Version    string
	Output     string
	MaxDepth   int
	Filter     string
}

// Parse получает конфигурацию из опций командной строки и проверяет её.
// При запросе справки выводит список опций и возвращает flag.ErrHelp.
func Parse(args []string) (Config, error) {
	var cfg Config
	flags := flag.NewFlagSet("dependency-visualizer", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&cfg.Package, "package", "", "имя анализируемого пакета (обязательно)")
	flags.StringVar(&cfg.Repository, "repository", "", "URL или путь к тестовому файлу (обязательно)")
	flags.BoolVar(&cfg.TestMode, "test", false, "использовать тестовый репозиторий")
	flags.StringVar(&cfg.Version, "version", "", "версия пакета (обязательно)")
	flags.StringVar(&cfg.Output, "output", "graph.png", "имя выходного PNG-файла")
	flags.IntVar(&cfg.MaxDepth, "max-depth", defaultMaxDepth, "максимальная глубина анализа")
	flags.StringVar(&cfg.Filter, "filter", "", "подстрока для исключения пакетов")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			flags.SetOutput(os.Stdout)
			flags.PrintDefaults()
		}
		return cfg, err
	}
	if flags.NArg() != noArguments {
		return cfg, fmt.Errorf("ожидались опции, получены аргументы: %v", flags.Args())
	}
	return cfg, cfg.Validate()
}

// Validate проверяет значения параметров без получения зависимостей пакета.
func (cfg Config) Validate() error {
	if !validToken(cfg.Package) {
		return fmt.Errorf("-package: требуется непустое имя без пробелов")
	}
	if err := cfg.validateRepository(); err != nil {
		return err
	}
	if !validToken(cfg.Version) {
		return fmt.Errorf("-version: требуется непустая версия без пробелов")
	}
	if filepath.Ext(cfg.Output) != ".png" || filepath.Base(cfg.Output) == ".png" {
		return fmt.Errorf("-output: требуется имя файла с расширением .png")
	}
	if cfg.MaxDepth < minDepth {
		return fmt.Errorf("-max-depth: глубина не может быть отрицательной")
	}
	if cfg.Filter != "" && !validToken(cfg.Filter) {
		return fmt.Errorf("-filter: подстрока не должна содержать пробелы")
	}
	return nil
}

// validateRepository проверяет URL или наличие обычного тестового файла.
func (cfg Config) validateRepository() error {
	if strings.TrimSpace(cfg.Repository) == "" {
		return fmt.Errorf("-repository: требуется URL или путь к файлу")
	}
	if cfg.TestMode {
		info, err := os.Stat(cfg.Repository)
		if err != nil {
			return fmt.Errorf("-repository: %w", err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("-repository: путь должен указывать на обычный файл")
		}
		return nil
	}
	address, err := url.ParseRequestURI(cfg.Repository)
	if err != nil {
		return fmt.Errorf("-repository: неверный URL: %w", err)
	}
	if (address.Scheme != "http" && address.Scheme != "https") || address.Hostname() == "" {
		return fmt.Errorf("-repository: требуется абсолютный HTTP(S) URL")
	}
	return nil
}

// validToken проверяет, что значение непустое и не содержит пробельных символов.
func validToken(value string) bool {
	return value != "" && !strings.ContainsAny(value, " \t\r\n")
}
