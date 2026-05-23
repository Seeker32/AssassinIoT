package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"ariga.io/atlas/atlasexec"
	"github.com/spf13/cobra"
)

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show which migrations are applied vs pending",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAtlasClient()
		if err != nil {
			return err
		}

		absDir, err := absMigrateDir()
		if err != nil {
			return err
		}

		params := &atlasexec.MigrateStatusParams{
			Env:    env,
			DirURL: "file://" + absDir,
		}
		if databaseURL != "" {
			params.URL = databaseURL
		}

		result, err := client.MigrateStatus(cmd.Context(), params)
		if err != nil {
			return fmt.Errorf("migration status check failed: %w", err)
		}

		if result.Error != "" {
			fmt.Fprintf(os.Stderr, "Warning: %s\n", result.Error)
			if result.SQL != "" {
				fmt.Fprintf(os.Stderr, "SQL: %s\n", result.SQL)
			}
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintf(w, "Version\tDescription\tStatus\tApplied At\n")
		fmt.Fprintf(w, "-------\t-----------\t------\t----------\n")
		for _, rev := range result.Applied {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				rev.Version, rev.Description, "APPLIED", rev.ExecutedAt.Format("2006-01-02 15:04:05"))
		}
		for _, f := range result.Pending {
			fmt.Fprintf(w, "%s\t%s\tPENDING\t-\n", f.Version, f.Description)
		}
		w.Flush()

		fmt.Printf("\nStatus: %s | Current: %s | Next: %s | Total: %d\n",
			result.Status, result.Current, result.Next, result.Total)
		return nil
	},
}

func init() {
	migrateCmd.AddCommand(migrateStatusCmd)
}
