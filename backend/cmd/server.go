package cmd

import (
	"github.com/Seeker32/AssassinIoT/backend/internal/application"
	"github.com/Seeker32/AssassinIoT/backend/internal/dependence"
	"github.com/spf13/cobra"
)

const (
	defaultServerConfigPath = "config.yaml"
)

var serverConfigPath string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the backend HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := resolveServerConfigPath(serverConfigPath)
		dep := dependence.NewDependence(dependence.WithConfigPath(cfgPath))
		appServer := application.NewServer(dep)
		return appServer.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVar(&serverConfigPath, "config", defaultServerConfigPath, "Configuration file path")
}

func resolveServerConfigPath(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return defaultServerConfigPath
}
