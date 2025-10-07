package cmds

import (
	"fmt"
	"os"
	"path/filepath"

	"patch-review-go/internal"
)

func CmdApply(configDir, reviewPath, patchPath string) {
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

	// Check if patch file exists
	if _, err := os.Stat(patchPath); os.IsNotExist(err) {
		fmt.Printf("Error: Patch file not found: %s\n", patchPath)
		return
	}

	// Detect patch format
	patchFormat := internal.DetectPatchFormat(patchPath)
	fmt.Printf("Detected patch format: %s\n", patchFormat.String())
	fmt.Printf("Applying patch: %s\n", patchPath)

	var applyErr error
	switch patchFormat {
	case internal.PatchFormatMbox:
		applyErr = internal.GitAmPatch(canonicalReviewPath, patchPath)
	case internal.PatchFormatDiff:
		applyErr = internal.GitApplyPatch(canonicalReviewPath, patchPath)
	default:
		fmt.Printf("Error: Unknown patch format: %s\n", patchPath)
		return
	}

	if applyErr != nil {
		fmt.Printf("Patch application completed with errors: %v\n", applyErr)
		fmt.Println("If there were conflicts, resolve them and use git commands to continue")
	} else {
		fmt.Println("Patch applied successfully")
	}
}