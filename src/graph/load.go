package graph

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Load открывает текстовый файл тестового репозитория и читает его граф.
func Load(path string) (Graph, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("тестовый репозиторий: %w", err)
	}
	defer file.Close()
	return Read(file)
}

// Read разбирает строки вида A: B C с именами из заглавных латинских букв.
func Read(reader io.Reader) (Graph, error) {
	source := make(Graph)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		name, value, ok := strings.Cut(line, ":")
		name = strings.TrimSpace(name)
		if !ok || !validName(name) {
			return nil, fmt.Errorf("неверная запись тестового графа: %q", line)
		}
		if _, exists := source[name]; exists {
			return nil, fmt.Errorf("повторное описание пакета %s", name)
		}
		dependencies := strings.Fields(value)
		for _, dependency := range dependencies {
			if !validName(dependency) {
				return nil, fmt.Errorf("неверное имя зависимости: %q", dependency)
			}
		}
		source[name] = dependencies
	}
	return source, scanner.Err()
}

// validName проверяет непустое имя из символов A-Z.
func validName(name string) bool {
	if name == "" {
		return false
	}
	for _, letter := range name {
		if letter < 'A' || letter > 'Z' {
			return false
		}
	}
	return true
}
