package cmd

import "github.com/spf13/cobra"

var migrateCmd = &cobra.Command{
	Use:     "migrate",
	Aliases: []string{"m"},
	Short:   "Manage database schema migrations",
	Long: `Manage database schema migrations using Atlas.

Supports status checks, applying pending migrations,
generating diffs from ent schema changes, and validating
migration directory integrity.`,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
