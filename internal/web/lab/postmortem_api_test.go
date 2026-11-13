package lab

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"breaklab/internal/progress"
	"breaklab/internal/web/webcore"
)

func postmortemAPILabYAML(id string) string {
	return "id: " + id + `
title: The website returns nothing
stage: linux
difficulty: 1
estimated_minutes: 5
environment:
  dockerfile: Dockerfile
  ports: []
check:
  script: check.sh
  run_in: host
`
}

// postmortemAPITestConfig builds a Config whose LabsDir/RoadmapPath are
// real files on disk — runner.SubmitPostmortem loads the catalog itself
// (via findLab), unlike buildLabDetail elsewhere in this package, which
// takes an already-loaded *catalog.Catalog.
func postmortemAPITestConfig(t *testing.T) (cfg webcore.Config, progressPath string) {
	t.Helper()
	root := t.TempDir()
	labDir := filepath.Join(root, "labs", "web-down")
	if err := os.MkdirAll(labDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(labDir, "lab.yaml"), []byte(postmortemAPILabYAML("web-down")), 0o644); err != nil {
		t.Fatal(err)
	}
	roadmapYAML := "phases:\n  - id: phase-1\n    title: Phase 1\nstages:\n  - id: linux\n    phase: phase-1\n    title: Linux fundamentals\n"
	if err := os.WriteFile(filepath.Join(root, "roadmap.yaml"), []byte(roadmapYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	progressPath = filepath.Join(root, "progress", "progress.json")
	cfg = webcore.Config{
		LabsDir:        filepath.Join(root, "labs"),
		RoadmapPath:    filepath.Join(root, "roadmap.yaml"),
		ProgressPath:   progressPath,
		PostmortemsDir: filepath.Join(root, "postmortems"),
	}
	return cfg, progressPath
}

func seedSolvedAttempt(t *testing.T, progressPath string) {
	t.Helper()
	solved := time.Now().UTC()
	pf := &progress.ProgressFile{Version: 1, Attempts: []progress.Attempt{
		{Lab: "web-down", Status: progress.StatusSolved, StartedAt: solved.Add(-10 * time.Minute), SolvedAt: &solved, HintsUsed: 1},
	}}
	if err := progress.Save(progressPath, pf); err != nil {
		t.Fatal(err)
	}
}

// validPostmortemSections fills every postmortem.RequiredSections key
// with enough words (ValidateSections requires >= 20 per section) to
// pass validation.
func validPostmortemSections() map[string]string {
	filler := "This is a sufficiently long section with more than twenty words in " +
		"total so it satisfies the postmortem validator's minimum length requirement for this field."
	return map[string]string{
		"Summary":    filler,
		"Detection":  filler,
		"Root cause": filler,
		"Fix":        filler,
		"Prevention": filler,
	}
}

func postPostmortem(cfg webcore.Config, labID string, sections map[string]string) *httptest.ResponseRecorder {
	body, _ := json.Marshal(postmortemSubmitRequest{Sections: sections})
	req := httptest.NewRequest("POST", "/api/labs/"+labID+"/postmortem", bytes.NewReader(body))
	req.SetPathValue("id", labID)
	rec := httptest.NewRecorder()
	labPostmortemAPIHandler(cfg)(rec, req)
	return rec
}

func TestLabPostmortemAPISuccessfulSubmissionMarksComplete(t *testing.T) {
	cfg, progressPath := postmortemAPITestConfig(t)
	seedSolvedAttempt(t, progressPath)

	rec := postPostmortem(cfg, "web-down", validPostmortemSections())
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body)
	}
	var resp postmortemSubmitResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !resp.Valid || len(resp.Failing) != 0 {
		t.Fatalf("response = %+v, want Valid with no failing sections", resp)
	}

	pf, err := progress.Load(progressPath)
	if err != nil {
		t.Fatal(err)
	}
	if pf.Attempts[0].Status != progress.StatusComplete {
		t.Errorf("attempt status = %q, want %q (a valid submission promotes solved -> complete)", pf.Attempts[0].Status, progress.StatusComplete)
	}
	if pf.Attempts[0].Postmortem == "" {
		t.Fatal("attempt should record the written postmortem's path")
	}
	if _, err := os.Stat(pf.Attempts[0].Postmortem); err != nil {
		t.Errorf("postmortem file should exist on disk: %v", err)
	}
}

func TestLabPostmortemAPIIncompleteSectionReturnsFailingWithoutPromoting(t *testing.T) {
	cfg, progressPath := postmortemAPITestConfig(t)
	seedSolvedAttempt(t, progressPath)

	sections := validPostmortemSections()
	sections["Root cause"] = "too short"
	rec := postPostmortem(cfg, "web-down", sections)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 even when the submission is invalid (body: %s)", rec.Code, rec.Body)
	}
	var resp postmortemSubmitResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Valid {
		t.Error("response.Valid = true, want false for a too-short required section")
	}
	if !slices.Contains(resp.Failing, "Root cause") {
		t.Errorf("Failing = %v, want it to name Root cause", resp.Failing)
	}

	pf, err := progress.Load(progressPath)
	if err != nil {
		t.Fatal(err)
	}
	if pf.Attempts[0].Status != progress.StatusSolved {
		t.Errorf("attempt status = %q, want still %q (never promoted on a failing submission)", pf.Attempts[0].Status, progress.StatusSolved)
	}
}

func TestLabPostmortemAPIRejectsUnsolvedLab(t *testing.T) {
	cfg, _ := postmortemAPITestConfig(t)
	// No attempt seeded at all — never started, let alone solved.

	rec := postPostmortem(cfg, "web-down", validPostmortemSections())
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (body: %s)", rec.Code, rec.Body)
	}
}

func TestLabPostmortemAPIResubmissionOverwritesSameFile(t *testing.T) {
	cfg, progressPath := postmortemAPITestConfig(t)
	seedSolvedAttempt(t, progressPath)

	first := postPostmortem(cfg, "web-down", validPostmortemSections())
	if first.Code != http.StatusOK {
		t.Fatalf("first submission status = %d, want 200 (body: %s)", first.Code, first.Body)
	}
	pf, err := progress.Load(progressPath)
	if err != nil {
		t.Fatal(err)
	}
	firstPath := pf.Attempts[0].Postmortem

	sections := validPostmortemSections()
	sections["Fix"] += " Updated with more detail on the second pass."
	second := postPostmortem(cfg, "web-down", sections)
	if second.Code != http.StatusOK {
		t.Fatalf("second submission status = %d, want 200 (body: %s)", second.Code, second.Body)
	}
	pf, err = progress.Load(progressPath)
	if err != nil {
		t.Fatal(err)
	}
	if pf.Attempts[0].Postmortem != firstPath {
		t.Errorf("resubmission path = %q, want it to overwrite the same file %q, not create a new one", pf.Attempts[0].Postmortem, firstPath)
	}
}
