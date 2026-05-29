package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"ariga.io/atlas/atlasexec"
	"github.com/spf13/cobra"
)

var (
	databaseURL       string
	migrateDir        string
	migrateConfigFile string
	env               string
)

var rootCmd = &cobra.Command{
	Use:   "assassin",
	Short: "AssassinIoT backend management CLI",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
}

func newAtlasClient() (*atlasexec.Client, error) {
	return atlasexec.NewClient(".", "atlas")
}

func absMigrateDir() (string, error) {
	abs, err := filepath.Abs(migrateDir)
	if err != nil {
		return "", fmt.Errorf("resolving migration directory: %w", err)
	}
	if _, err := os.Stat(abs); err != nil {
		return "", fmt.Errorf("migration directory %s: %w", abs, err)
	}
	return abs, nil
}

func projectRoot() (string, error) {
	abs, err := absMigrateDir()
	if err != nil {
		return "", err
	}
	// Migration dir is typically <project>/ent/migrate/migrations — go up 3 levels.
	// Also handle the case where --dir is set to a custom path by walking up to find go.mod.
	dir := abs
	for range 4 {
		dir = filepath.Dir(dir)
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
	}
	return filepath.Dir(filepath.Dir(filepath.Dir(abs))), nil
}
