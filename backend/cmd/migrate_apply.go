package cmd

import (
	"fmt"
	"os"

	"ariga.io/atlas/atlasexec"
	"github.com/spf13/cobra"
)

var (
	applyAmount     uint64
	applyDryRun     bool
	applyAllowDirty bool
)

var migrateApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply pending migrations to the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAtlasClient()
		if err != nil {
			return err
		}

		absDir, err := absMigrateDir()
		if err != nil {
			return err
		}

		params := &atlasexec.MigrateApplyParams{
			Env:        env,
			DirURL:     "file://" + absDir,
			Amount:     applyAmount,
			DryRun:     applyDryRun,
			AllowDirty: applyAllowDirty,
		}
		if databaseURL != "" {
			params.URL = databaseURL
		}

		result, err := client.MigrateApply(cmd.Context(), params)
		if err != nil {
			return fmt.Errorf("migration apply failed: %w", err)
		}

		for _, f := range result.Applied {
			if f.Error != nil {
				fmt.Fprintf(os.Stderr, "ERROR applying %s: %s\n", f.Name, f.Error.Text)
			} else {
				dur := f.End.Sub(f.Start)
				if applyDryRun {
					fmt.Printf("Would apply %s\n", f.Name)
					for _, stmt := range f.Applied {
						fmt.Println(stmt)
					}
				} else {
					fmt.Printf("Applied %s (%v)\n", f.Name, dur)
				}
			}
		}

		if result.Error != "" {
			return fmt.Errorf("apply completed with errors: %s", result.Error)
		}
		if len(result.Applied) == 0 {
			fmt.Println("No pending migrations to apply.")
		}
		return nil
	},
}

func init() {
	migrateCmd.AddCommand(migrateApplyCmd)
	migrateApplyCmd.Flags().Uint64Var(&applyAmount, "amount", 0,
		"Number of migrations to apply (0 = all pending)")
	migrateApplyCmd.Flags().BoolVar(&applyDryRun, "dry-run", false,
		"Print SQL without applying to the database")
	migrateApplyCmd.Flags().BoolVar(&applyAllowDirty, "allow-dirty", false,
		"Allow applying on a database with a dirty migration state")
}
