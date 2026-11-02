package apt

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ReadGraph читает индекс, выбирая точную версию корня и первую запись остальных имён.
func ReadGraph(reader io.Reader, root, version string) (map[string][]string, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(nil, maxLineSize)
	groups := make(map[string][]string)
	providers := make(map[string]string)
	found := false
	for {
		fields, err := readFields(scanner)
		if err != nil {
			return nil, fmt.Errorf("чтение индекса: %w", err)
		}
		if len(fields) == noFields {
			break
		}
		name := fields["package"]
		if name == "" || fields["version"] == "" {
			return nil, fmt.Errorf("запись индекса без Package или Version")
		}
		recordProvides(fields, providers)
		if _, exists := groups[name]; !exists {
			groups[name] = parseDependencies(fields)
		}
		if name == root && fields["version"] == version {
			groups[name] = parseDependencies(fields)
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("пакет %s версии %s не найден в индексе", root, version)
	}
	return resolveNames(groups, providers), nil
}

// resolveNames извлекает имена, выбирая первую имеющуюся в индексе альтернативу.
func resolveNames(groups map[string][]string, providers map[string]string) map[string][]string {
	result := make(map[string][]string)
	for name, dependencies := range groups {
		result[name] = nil
		for _, group := range dependencies {
			chosen := ""
			for _, alternative := range strings.Split(group, "|") {
				candidate := dependencyName(alternative)
				if candidate == "" {
					continue
				}
				if chosen == "" {
					chosen = candidate
				}
				if _, exists := groups[candidate]; exists {
					chosen = candidate
					break
				}
				if provider := providers[candidate]; provider != "" {
					chosen = provider
					break
				}
			}
			if chosen != "" {
				result[name] = append(result[name], chosen)
			}
		}
	}
	return result
}

// recordProvides запоминает первый пакет-поставщик каждого виртуального имени.
func recordProvides(fields map[string]string, providers map[string]string) {
	for _, value := range strings.Split(fields["provides"], ",") {
		name := dependencyName(value)
		if name != "" && providers[name] == "" {
			providers[name] = fields["package"]
		}
	}
}

// dependencyName извлекает имя без ограничения версии и квалификатора архитектуры.
func dependencyName(value string) string {
	text, _, _ := strings.Cut(value, "(")
	fields := strings.Fields(text)
	if len(fields) == noFields {
		return ""
	}
	name, _, _ := strings.Cut(fields[0], ":")
	return name
}
