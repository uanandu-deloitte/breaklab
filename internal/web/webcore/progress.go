package webcore

import (
	"breaklab/internal/postmortem"
	"breaklab/internal/progress"
)

// LoadProgress loads progress.json for a request and confines every
// recorded postmortem path to cfg.PostmortemsDir: a valid one is replaced
// with its confined absolute path, anything else (outside the directory,
// a symlink, a missing file) is blanked. Every handler loads progress
// through here, so no view can link to or read a path that a
// hand-edited progress.json points elsewhere.
func LoadProgress(cfg Config) (*progress.ProgressFile, error) {
	pf, err := progress.Load(cfg.ProgressPath)
	if err != nil {
		return nil, err
	}
	for i := range pf.Attempts {
		recorded := pf.Attempts[i].Postmortem
		if recorded == "" {
			continue
		}
		confined, err := postmortem.Confined(cfg.PostmortemsDir, recorded)
		if err != nil {
			confined = ""
		}
		pf.Attempts[i].Postmortem = confined
	}
	return pf, nil
}
