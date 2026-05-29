package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Seeker32/AssassinIoT/backend/internal/conf"
	"github.com/Seeker32/AssassinIoT/backend/internal/dependence"
	"github.com/spf13/cobra"
)

const (
	defaultServerAddr       = ":8080"
	defaultServerConfigPath = "config.yaml"
	shutdownTimeout         = 5 * time.Second
)

var serverAddr string
var serverConfigPath string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the backend HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := resolveServerConfigPath(serverConfigPath)
		dep := dependence.NewDependence(dependence.WithConfigPath(cfgPath))

		cfg := dep.ConfigProvider().ServerConfig()
		addr := resolveServerAddr(serverAddr, cfg)
		logger := dep.Logger()

		mux := http.NewServeMux()
		srv := &http.Server{
			Addr:    addr,
			Handler: mux,
		}

		errCh := make(chan error, 1)
		go func() {
			logger.Info("starting server", "addr", addr, "config", cfgPath)
			errCh <- srv.ListenAndServe()
		}()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigCh)

		select {
		case err := <-errCh:
			if err != nil && err != http.ErrServerClosed {
				return fmt.Errorf("start server: %w", err)
			}
			return nil
		case sig := <-sigCh:
			logger.Info("shutting down server", "signal", sig.String())
			ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				return fmt.Errorf("shutdown server: %w", err)
			}
			if err := <-errCh; err != nil && err != http.ErrServerClosed {
				return fmt.Errorf("server exited after shutdown: %w", err)
			}
			return nil
		case <-cmd.Context().Done():
			ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				return fmt.Errorf("shutdown server: %w", err)
			}
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVar(&serverAddr, "addr", "", "HTTP listen address")
	serverCmd.Flags().StringVar(&serverConfigPath, "config", defaultServerConfigPath, "Configuration file path")
}

func resolveServerConfigPath(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return defaultServerConfigPath
}

func resolveServerAddr(flagValue string, cfg conf.ServerConfig) string {
	if flagValue != "" {
		return flagValue
	}
	if cfg.Addr != "" {
		return cfg.Addr
	}
	return defaultServerAddr
}
