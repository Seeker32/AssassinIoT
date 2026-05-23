package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"ariga.io/atlas/atlasexec"
	"github.com/spf13/cobra"
)

var (
	diffDevURL string
	diffToURL  string
)

var migrateDiffCmd = &cobra.Command{
	Use:   "diff <name>",
	Short: "Generate a new migration by comparing ent schema with the database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		absDir, err := absMigrateDir()
		if err != nil {
			return err
		}

		pr, err := projectRoot()
		if err != nil {
			return err
		}

		// Use exec.Command instead of SDK MigrateDiff because SDK requires cloud login.
		// The atlas CLI binary handles local diff generation (file writing) directly.
		atlasArgs := []string{
			"migrate", "diff",
			"--env", env,
			"--dir", "file://" + absDir,
		}
		if diffDevURL != "" {
			atlasArgs = append(atlasArgs, "--dev-url", diffDevURL)
		}
		if diffToURL != "" {
			atlasArgs = append(atlasArgs, "--to", diffToURL)
		}
		atlasArgs = append(atlasArgs, name)

		atlasCmd := exec.CommandContext(cmd.Context(), "atlas", atlasArgs...)
		atlasCmd.Dir = pr
		atlasCmd.Stdout = os.Stdout
		atlasCmd.Stderr = os.Stderr

		if err := atlasCmd.Run(); err != nil {
			return fmt.Errorf("atlas migrate diff failed: %w", err)
		}

		// Recompute atlas.sum after diff generates new files.
		client, err := newAtlasClient()
		if err != nil {
			return err
		}
		if err := client.MigrateHash(cmd.Context(), &atlasexec.MigrateHashParams{
			DirURL: "file://" + absDir,
		}); err != nil {
			return fmt.Errorf("updating atlas.sum failed: %w", err)
		}
		fmt.Println("Updated atlas.sum")
		return nil
	},
}

func init() {
	migrateCmd.AddCommand(migrateDiffCmd)
	migrateDiffCmd.Flags().StringVar(&diffDevURL, "dev-url", "",
		"Dev database URL (default: from atlas.hcl environment)")
	migrateDiffCmd.Flags().StringVar(&diffToURL, "to", "ent://ent/schema",
		"Desired schema URL")
}
