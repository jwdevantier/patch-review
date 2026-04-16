package cmds

import (
	"fmt"
	"path/filepath"

	"github.com/jwdevantier/patch-review/internal"
)

func CmdReset(configDir, reviewPath, sourceAlias string) {
	config, err := internal.LoadConfig(configDir)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	state, err := internal.LoadState(configDir)
	if err != nil {
		fmt.Printf("Error loading state: %v\n", err)
		return
	}

	// Use default source if none specified
	if sourceAlias == "" {
		sourceAlias = config.GetDefaultSource()
		if sourceAlias == "" {
			fmt.Println("Error: No source specified and no default source configured")
			return
		}
	}

	source, err := config.GetSource(sourceAlias)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
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

	repoPath := internal.ExpandPathString(source.Path)
	baseBranch := source.Branch
	baseRemote := source.Remote
	branchPrefix := config.GetBranchPrefix()
	branchName := internal.GenerateBranchName(branchPrefix)

	// Check for existing worktree
	existingWorktree := state.GetWorktree(canonicalReviewPath)
	if existingWorktree != nil {
		fmt.Printf("Removing existing worktree: %s\n", canonicalReviewPath)
		if err := internal.GitWorktreeRemove(repoPath, canonicalReviewPath); err != nil {
			fmt.Printf("Warning: Failed to remove worktree: %v\n", err)
		}
		if err := internal.GitDeleteBranch(repoPath, existingWorktree.Branch); err != nil {
			fmt.Printf("Warning: Failed to delete branch: %v\n", err)
		}
	}

	// Fetch latest changes
	fmt.Println("Fetching latest changes...")
	if err := internal.GitFetch(repoPath, baseRemote); err != nil {
		fmt.Printf("Error fetching: %v\n", err)
		return
	}

	// Create new branch
	fmt.Printf("Creating branch: %s\n", branchName)
	if err := internal.GitCreateBranch(repoPath, branchName, baseRemote, baseBranch); err != nil {
		fmt.Printf("Error creating branch: %v\n", err)
		return
	}

	// Create worktree
	fmt.Printf("Creating worktree: %s\n", canonicalReviewPath)
	if err := internal.GitWorktreeAdd(repoPath, canonicalReviewPath, branchName); err != nil {
		fmt.Printf("Error creating worktree: %v\n", err)
		return
	}

	// Update state
	fmt.Println("Updating state...")
	worktree := internal.MakeWorktree(canonicalReviewPath, repoPath, branchName, baseBranch, sourceAlias)
	state.AddWorktree(worktree)
	if err := state.SaveState(configDir); err != nil {
		fmt.Printf("Error saving state: %v\n", err)
		return
	}

	fmt.Printf("Review directory ready: %s\n", canonicalReviewPath)
}
