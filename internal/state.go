package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Worktree struct {
	Path         string `json:"path"`
	Repo         string `json:"repo"`
	Branch       string `json:"branch"`
	BaseBranch   string `json:"base-branch"`
	SourceAlias  string `json:"source-alias"`
	Created      string `json:"created"`
}

type State struct {
	Worktrees []Worktree `json:"worktrees"`
}

func LoadState(configDir string) (*State, error) {
	statePath := filepath.Join(configDir, "patch-review.state.json")

	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		return &State{Worktrees: []Worktree{}}, nil
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %v", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %v", err)
	}

	return &state, nil
}

func (s *State) SaveState(configDir string) error {
	statePath := filepath.Join(configDir, "patch-review.state.json")

	if err := os.MkdirAll(filepath.Dir(statePath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %v", err)
	}

	if err := os.WriteFile(statePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %v", err)
	}

	return nil
}

func (s *State) AddWorktree(wt Worktree) {
	// Remove existing worktree with same path if it exists
	s.RemoveWorktree(wt.Path)
	s.Worktrees = append(s.Worktrees, wt)
}

func (s *State) RemoveWorktree(path string) {
	for i, wt := range s.Worktrees {
		if wt.Path == path {
			s.Worktrees = append(s.Worktrees[:i], s.Worktrees[i+1:]...)
			return
		}
	}
}

func (s *State) GetWorktree(path string) *Worktree {
	for _, wt := range s.Worktrees {
		if wt.Path == path {
			return &wt
		}
	}
	return nil
}

func MakeWorktree(path, repo, branch, baseBranch, sourceAlias string) Worktree {
	// Clean and canonicalize paths
	cleanPath := filepath.Clean(path)
	if !filepath.IsAbs(cleanPath) {
		abs, _ := filepath.Abs(cleanPath)
		cleanPath = abs
	}
	if !strings.HasSuffix(cleanPath, string(filepath.Separator)) {
		cleanPath += string(filepath.Separator)
	}

	cleanRepo := filepath.Clean(ExpandPathString(repo))
	if !filepath.IsAbs(cleanRepo) {
		abs, _ := filepath.Abs(cleanRepo)
		cleanRepo = abs
	}
	if !strings.HasSuffix(cleanRepo, string(filepath.Separator)) {
		cleanRepo += string(filepath.Separator)
	}

	return Worktree{
		Path:        cleanPath,
		Repo:        cleanRepo,
		Branch:      branch,
		BaseBranch:  baseBranch,
		SourceAlias: sourceAlias,
		Created:     time.Now().Format("Monday, January 2nd, 2006 3:04:05pm"),
	}
}