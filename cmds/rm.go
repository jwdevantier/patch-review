package cmds

import (
	"fmt"
	"path/filepath"

	"patch-review-go/internal"
)

func CmdRm(configDir, reviewPath string) {
	state, err := internal.LoadState(configDir)
	if err != nil {
		fmt.Printf("Error loading state: %v\n", err)
		return
	}

	// Canonicalize review path
	absReviewPath, err := filepath.Abs(reviewPath)
	if err != nil {
		fmt.Printf("Error getting absolute path: %v\n", err)
		return
	}
	canonicalReviewPath := filepath.Clean(absReviewPath)
	if !filepath.IsAbs(canonicalReviewPath) {
		fmt.Printf("Error: Path must be absolute\n")
		return
	}
	if canonicalReviewPath[len(canonicalReviewPath)-1] != filepath.Separator {
		canonicalReviewPath += string(filepath.Separator)
	}

	worktree := state.GetWorktree(canonicalReviewPath)
	if worktree == nil {
		fmt.Printf("Error: Review directory not found in state: %s\n", canonicalReviewPath)
		return
	}

	repoPath := worktree.Repo
	branchName := worktree.Branch

	// Remove worktree
	fmt.Printf("Removing worktree: %s\n", canonicalReviewPath)
	if err := internal.GitWorktreeRemove(repoPath, canonicalReviewPath); err != nil {
		fmt.Printf("Error removing worktree: %v\n", err)
		return
	}

	// Delete branch
	fmt.Printf("Deleting branch: %s\n", branchName)
	if err := internal.GitDeleteBranch(repoPath, branchName); err != nil {
		fmt.Printf("Error deleting branch: %v\n", err)
		return
	}

	// Update state
	state.RemoveWorktree(canonicalReviewPath)
	if err := state.SaveState(configDir); err != nil {
		fmt.Printf("Error saving state: %v\n", err)
		return
	}

	fmt.Printf("Review directory removed: %s\n", canonicalReviewPath)
}