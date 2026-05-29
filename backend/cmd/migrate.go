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
	migrateCmd.PersistentFlags().StringVarP(&databaseURL, "url", "u", "",
		"Database connection URL (overrides atlas.hcl environment)")
	migrateCmd.PersistentFlags().StringVar(&migrateDir, "dir", "ent/migrate/migrations",
		"Migration directory path")
	migrateCmd.PersistentFlags().StringVar(&migrateConfigFile, "config", "ent/migrate/atlas.hcl",
		"Atlas HCL config file path")
	migrateCmd.PersistentFlags().StringVar(&env, "env", "local",
		"Atlas environment name from HCL config")
}
