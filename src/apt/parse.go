package apt

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	maxLineSize = 1024 * 1024
	noFields    = 0
)

// Find ищет пакет по точному имени и версии и возвращает группы прямых зависимостей.
func Find(reader io.Reader, name, version string) ([]string, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(nil, maxLineSize)
	for {
		fields, err := readFields(scanner)
		if err != nil {
			return nil, fmt.Errorf("чтение индекса: %w", err)
		}
		if len(fields) == noFields {
			return nil, fmt.Errorf("пакет %s версии %s не найден в индексе", name, version)
		}
		if fields["package"] == name && fields["version"] == version {
			return parseDependencies(fields), nil
		}
	}
}

// readFields читает одну запись индекса с учётом строк продолжения полей.
func readFields(scanner *bufio.Scanner) (map[string]string, error) {
	fields := make(map[string]string)
	key := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			if len(fields) == noFields {
				continue
			}
			return fields, nil
		}
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			if key == "" {
				return nil, fmt.Errorf("строка продолжения без поля")
			}
			fields[key] += " " + strings.TrimSpace(line)
			continue
		}
		name, value, ok := strings.Cut(line, ":")
		if !ok || name == "" {
			return nil, fmt.Errorf("неверная строка: %q", line)
		}
		key = strings.ToLower(name)
		fields[key] = strings.TrimSpace(value)
	}
	return fields, scanner.Err()
}

// parseDependencies объединяет Pre-Depends и Depends, сохраняя версии и альтернативы.
func parseDependencies(fields map[string]string) []string {
	var dependencies []string
	for _, key := range []string{"pre-depends", "depends"} {
		for _, value := range strings.Split(fields[key], ",") {
			value = strings.TrimSpace(value)
			if value != "" {
				dependencies = append(dependencies, value)
			}
		}
	}
	return dependencies
}
