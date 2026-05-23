package cmd

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

type hashEntry struct{ N, H string }

var migrateValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Verify migration directory integrity (check atlas.sum hashes)",
	RunE: func(cmd *cobra.Command, args []string) error {
		absDir, err := absMigrateDir()
		if err != nil {
			return err
		}

		// Parse atlas.sum.
		sumPath := filepath.Join(absDir, "atlas.sum")
		data, err := os.ReadFile(sumPath)
		if err != nil {
			return fmt.Errorf("reading %s: %w", sumPath, err)
		}

		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		if len(lines) < 2 {
			return fmt.Errorf("%s: must have at least 2 lines (directory hash + file entries)", sumPath)
		}

		storedDirHash := strings.TrimPrefix(lines[0], "h1:")
		var storedEntries []hashEntry
		for _, line := range lines[1:] {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, " h1:", 2)
			if len(parts) != 2 {
				fmt.Fprintf(os.Stderr, "WARN: malformed line: %s\n", line)
				continue
			}
			storedEntries = append(storedEntries, hashEntry{strings.TrimSpace(parts[0]), parts[1]})
		}

		// Collect .sql files sorted by name.
		entries, err := os.ReadDir(absDir)
		if err != nil {
			return fmt.Errorf("reading migration directory: %w", err)
		}
		var sqlFiles []string
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
				sqlFiles = append(sqlFiles, e.Name())
			}
		}
		sort.Strings(sqlFiles)

		// Compute cumulative hashes — matches Atlas NewHashFile algorithm.
		var computedEntries []hashEntry
		var errorsFound int
		h := sha256.New()

		for _, name := range sqlFiles {
			content, err := os.ReadFile(filepath.Join(absDir, name))
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: cannot read %s: %v\n", name, err)
				errorsFound++
				continue
			}
			h.Write([]byte(name))
			h.Write(content)
			computedEntries = append(computedEntries, hashEntry{
				N: name,
				H: base64.StdEncoding.EncodeToString(h.Sum(nil)),
			})
		}

		// Compare per-file hashes.
		storedByName := map[string]hashEntry{}
		for _, e := range storedEntries {
			storedByName[e.N] = e
		}
		for _, c := range computedEntries {
			s, ok := storedByName[c.N]
			if !ok {
				fmt.Fprintf(os.Stderr, "ERROR: %s: not in atlas.sum (file was added)\n", c.N)
				errorsFound++
				continue
			}
			if c.H != s.H {
				fmt.Fprintf(os.Stderr, "ERROR: %s: hash mismatch\nexpected:  %s\nactual:    %s\n",
					c.N, s.H, c.H)
				errorsFound++
			} else {
				fmt.Printf("OK  %s\n", c.N)
			}
		}
		for _, s := range storedEntries {
			found := false
			for _, c := range computedEntries {
				if c.N == s.N {
					found = true
					break
				}
			}
			if !found {
				fmt.Fprintf(os.Stderr, "ERROR: %s: in atlas.sum but missing (file was removed)\n", s.N)
				errorsFound++
			}
		}

		// Compute and verify directory hash.
		dirHash := sha256.New()
		for _, e := range computedEntries {
			dirHash.Write([]byte(e.N))
			dirHash.Write([]byte(e.H))
		}
		actualDirHash := base64.StdEncoding.EncodeToString(dirHash.Sum(nil))

		if actualDirHash != storedDirHash {
			fmt.Fprintf(os.Stderr, "ERROR: directory hash mismatch\nexpected:  %s\nactual:    %s\n",
				storedDirHash, actualDirHash)
			errorsFound++
		} else {
			fmt.Println("OK  directory hash")
		}

		if errorsFound > 0 {
			return fmt.Errorf("validation failed with %d error(s)", errorsFound)
		}
		fmt.Println("\nMigration directory integrity check PASSED.")
		return nil
	},
}

func init() {
	migrateCmd.AddCommand(migrateValidateCmd)
}
