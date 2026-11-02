package main

import (
	"bytes"

	"github.com/marseliandrius/dependency-visualizer/src/apt"
	"github.com/marseliandrius/dependency-visualizer/src/config"
	"github.com/marseliandrius/dependency-visualizer/src/graph"
)

// loadRepository читает тестовый граф или индекс Ubuntu для указанной версии пакета.
func loadRepository(cfg config.Config) (graph.Graph, error) {
	if cfg.TestMode {
		return graph.Load(cfg.Repository)
	}
	data, err := apt.Download(cfg.Repository)
	if err != nil {
		return nil, err
	}
	return apt.ReadGraph(bytes.NewReader(data), cfg.Package, cfg.Version)
}
