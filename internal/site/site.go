package site

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"breaklab/internal/catalog"
	"breaklab/internal/progress"
)

// Config locates everything Build needs to read and where it writes.
type Config struct {
	LabsDir      string
	RoadmapPath  string
	ProgressPath string
	TemplatesDir string
	OutDir       string

	// PostmortemsDir confines which files a progress.json postmortem
	// path may name; anything outside it is never read.
	PostmortemsDir string
}

// Build rebuilds OutDir from scratch (wiping any previous contents —
// design.md Decision 5, so a postmortem un-published after being
// published can't leave a stale page behind) using only RoadmapPath,
// lab metadata, ProgressPath, and postmortems marked publish: true.
func Build(cfg Config) error {
	if err := os.RemoveAll(cfg.OutDir); err != nil {
		return fmt.Errorf("clearing %s: %w", cfg.OutDir, err)
	}
	if err := os.MkdirAll(cfg.OutDir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", cfg.OutDir, err)
	}

	cat, err := catalog.LoadCatalog(cfg.LabsDir, cfg.RoadmapPath)
	if err != nil {
		return err
	}
	pf, err := progress.Load(cfg.ProgressPath)
	if err != nil {
		return err
	}

	published, views, err := LoadPublished(cfg.PostmortemsDir, CollectPostmortemPaths(pf))
	if err != nil {
		return err
	}

	if err := writeIndex(cfg, BuildIndex(cat, pf, published)); err != nil {
		return err
	}
	if err := writePostmortems(cfg, views); err != nil {
		return err
	}
	if err := writePostmortemsIndex(cfg, BuildPostmortemsIndex(cat, pf, published)); err != nil {
		return err
	}
	if err := copyStyle(cfg); err != nil {
		return err
	}
	return copyFonts(cfg)
}

func writeIndex(cfg Config, view IndexView) error {
	tmpl, err := template.ParseFiles(filepath.Join(cfg.TemplatesDir, "index.html"))
	if err != nil {
		return err
	}
	return renderToFile(tmpl, filepath.Join(cfg.OutDir, "index.html"), view)
}

func writePostmortems(cfg Config, views map[string]PostmortemView) error {
	if len(views) == 0 {
		return nil
	}
	tmpl, err := template.ParseFiles(filepath.Join(cfg.TemplatesDir, "postmortem.html"))
	if err != nil {
		return err
	}
	dir := filepath.Join(cfg.OutDir, "postmortems")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}
	for _, view := range views {
		path := filepath.Join(dir, view.Slug+".html")
		if err := renderToFile(tmpl, path, view); err != nil {
			return err
		}
	}
	return nil
}

// writePostmortemsIndex writes the flat, public postmortems index to
// site/postmortems/index.html, creating the postmortems/ directory
// itself (writePostmortems only creates it when there's at least one
// published postmortem, but the index page renders — with an empty
// list — regardless of count).
func writePostmortemsIndex(cfg Config, view PostmortemsIndexView) error {
	tmpl, err := template.ParseFiles(filepath.Join(cfg.TemplatesDir, "postmortems.html"))
	if err != nil {
		return err
	}
	dir := filepath.Join(cfg.OutDir, "postmortems")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}
	return renderToFile(tmpl, filepath.Join(dir, "index.html"), view)
}

func copyStyle(cfg Config) error {
	data, err := os.ReadFile(filepath.Join(cfg.TemplatesDir, "style.css"))
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(cfg.OutDir, "style.css"), data, 0o644)
}

// vendoredFontFiles are style.css's @font-face assets — copied
// alongside it so the published site keeps the offline/no-external-
// request guarantee the web-ui-design spec's "Locally-vendored
// typography" requirement extends to this surface too.
var vendoredFontFiles = []string{
	"JetBrainsMono-Regular.woff2", "JetBrainsMono-Medium.woff2",
	"Geist-Regular.woff2", "Geist-SemiBold.woff2", "Geist-Bold.woff2",
}

func copyFonts(cfg Config) error {
	outDir := filepath.Join(cfg.OutDir, "fonts")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", outDir, err)
	}
	for _, name := range vendoredFontFiles {
		data, err := os.ReadFile(filepath.Join(cfg.TemplatesDir, "fonts", name))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(outDir, name), data, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func renderToFile(tmpl *template.Template, path string, data any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(f, data); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
