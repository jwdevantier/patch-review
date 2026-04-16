package internal

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func GitCommand(repoPath string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git command failed: git %s\nError: %s", strings.Join(args, " "), string(output))
	}
	return nil
}

func GitFetch(repoPath, remote string) error {
	return GitCommand(repoPath, "fetch", remote)
}

func GitCreateBranch(repoPath, branchName, baseRemote, baseBranch string) error {
	var src string
	if baseRemote != "" {
		src = fmt.Sprintf("%s/%s", baseRemote, baseBranch)
	} else {
		src = baseBranch
	}
	return GitCommand(repoPath, "branch", branchName, src)
}

func GitDeleteBranch(repoPath, branchName string) error {
	return GitCommand(repoPath, "branch", "-D", branchName)
}

func GitWorktreeAdd(repoPath, worktreePath, branchName string) error {
	return GitCommand(repoPath, "worktree", "add", worktreePath, branchName)
}

func GitWorktreeRemove(repoPath, worktreePath string) error {
	return GitCommand(repoPath, "worktree", "remove", worktreePath)
}

func GitApplyPatch(worktreePath, patchFile string) error {
	cmd := exec.Command("git", "apply", patchFile)
	cmd.Dir = worktreePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(output))
		return err
	}
	fmt.Print(string(output))
	return nil
}

func GitAmPatch(worktreePath, patchFile string) error {
	cmd := exec.Command("git", "am", patchFile)
	cmd.Dir = worktreePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(output))
		return err
	}
	fmt.Print(string(output))
	return nil
}

func GenerateBranchName(prefix string) string {
	now := time.Now()
	timestamp := fmt.Sprintf("%04d%02d%02d%02d%02d%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
	randomSuffix := fmt.Sprintf("%x", rand.Intn(4096))
	return fmt.Sprintf("%s/%s-%s", prefix, timestamp, randomSuffix)
}

type PatchFormat int

const (
	PatchFormatUnknown PatchFormat = iota
	PatchFormatMbox
	PatchFormatDiff
)

func (pf PatchFormat) String() string {
	switch pf {
	case PatchFormatMbox:
		return "mbox"
	case PatchFormatDiff:
		return "diff"
	default:
		return "unknown"
	}
}

func DetectPatchFormat(filePath string) PatchFormat {
	ext := filepath.Ext(filePath)

	// Check by extension first
	switch ext {
	case ".mbox", ".mbx", ".eml":
		return PatchFormatMbox
	case ".patch", ".diff":
		return PatchFormatDiff
	}

	// Inspect file content
	file, err := os.Open(filePath)
	if err != nil {
		return PatchFormatUnknown
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	mboxPattern := regexp.MustCompile(`^(From |Date:|Subject:|Message-ID:)`)
	diffPattern := regexp.MustCompile(`^(diff |Index:|--- |\+\+\+ )`)

	for scanner.Scan() && lineCount < 10 {
		line := scanner.Text()
		if mboxPattern.MatchString(line) {
			return PatchFormatMbox
		}
		if diffPattern.MatchString(line) {
			return PatchFormatDiff
		}
		lineCount++
	}

	return PatchFormatUnknown
}
