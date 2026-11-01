// Package apt получает прямые зависимости из индекса пакетов Ubuntu.
package apt

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const requestTimeout = 30 * time.Second

// Fetch загружает индекс Packages или Packages.gz и ищет заданную версию пакета.
func Fetch(address, name, version string) ([]string, error) {
	client := http.Client{Timeout: requestTimeout}
	response, err := client.Get(address)
	if err != nil {
		return nil, fmt.Errorf("загрузка индекса: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("загрузка индекса: HTTP %s", response.Status)
	}
	var reader io.Reader = response.Body
	if strings.HasSuffix(response.Request.URL.Path, ".gz") && !response.Uncompressed {
		archive, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, fmt.Errorf("чтение gzip: %w", err)
		}
		defer archive.Close()
		reader = archive
	}
	return Find(reader, name, version)
}
