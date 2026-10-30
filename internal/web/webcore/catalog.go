package webcore

import "breaklab/internal/catalog"

// LoadCatalog loads the catalog from cfg's configured paths.
func LoadCatalog(cfg Config) (*catalog.Catalog, error) {
	return catalog.LoadCatalog(cfg.LabsDir, cfg.RoadmapPath)
}
