package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"patch-review/cmds"
)

var configDir string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "patch-review",
		Short: "A tool for managing patch reviews with git worktrees",
		Long:  "patch-review simplifies the management of git working trees for patch review workflows",
	}

	rootCmd.PersistentFlags().StringVar(&configDir, "dir", "", "Path to config files")

	resetCmd := &cobra.Command{
		Use:   "reset [--source <name>] <review-dir>",
		Short: "Create or refresh a git worktree directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			source, _ := cmd.Flags().GetString("source")
			cmds.CmdReset(getConfigDir(), args[0], source)
		},
	}
	resetCmd.Flags().String("source", "", "Name of source to use")

	rmCmd := &cobra.Command{
		Use:   "rm <review-dir>",
		Short: "Remove a worktree from disk and clean up git repo state",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmds.CmdRm(getConfigDir(), args[0])
		},
	}

	applyCmd := &cobra.Command{
		Use:   "apply <review-dir> <patch-file>",
		Short: "Apply a git patch file (or mailbox file) to worktree",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cmds.CmdApply(getConfigDir(), args[0], args[1])
		},
	}

	rootCmd.AddCommand(resetCmd, rmCmd, applyCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func getConfigDir() string {
	if configDir != "" {
		return configDir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Failed to get home directory: %v", err))
	}
	return filepath.Join(home, ".config", "patch-review")
}
