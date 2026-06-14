package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all entries in the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := openVault()
		if err != nil {
			return err
		}
		entries, err := svc.ListEntries()
		if err != nil {
			return err
		}
		// print each EntryMetadata (id, username, url) — no passwords, by design
		for _, e := range entries {
			fmt.Printf("[%d] Username: %s - URL: %s\n", e.ID, e.Username, e.URL)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}