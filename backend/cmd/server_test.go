package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestServerCommandRegistered(t *testing.T) {
	t.Helper()

	cmd, _, err := rootCmd.Find([]string{"server"})
	if err != nil {
		t.Fatalf("find server command: %v", err)
	}
	if cmd == nil {
		t.Fatal("expected server command to be registered")
	}
	if cmd.Name() != "server" {
		t.Fatalf("expected command name server, got %q", cmd.Name())
	}
	if flag := cmd.Flags().Lookup("addr"); flag != nil {
		t.Fatal("expected server command to not expose --addr")
	}
}

func TestResolveServerConfigPath(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("server:\n  log_level: info\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	tests := []struct {
		name string
		flag string
		want string
	}{
		{
			name: "flag overrides default",
			flag: configPath,
			want: configPath,
		},
		{
			name: "default config path",
			want: defaultServerConfigPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveServerConfigPath(tt.flag)
			if got != tt.want {
				t.Fatalf("resolveServerConfigPath(%q) = %q, want %q", tt.flag, got, tt.want)
			}
		})
	}
}
