package cli 

import (
	"fmt"
	"strconv"
	"github.com/spf13/cobra"
)

var (
	updateUsername string
	updateURL      string
	updateNotes    string
)

var updateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update fields of an existing entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid id %q: %w", args[0], err)
		}
		svc, err := openVault()
		if err != nil {
			return err
		}

		// start from the CURRENT entry so unspecified fields are preserved
		entry, err := svc.GetEntry(id)
		if err != nil {
			return err
		}

		if cmd.Flags().Changed("username") {
			entry.Username = updateUsername
		}
		if cmd.Flags().Changed("url") {
			entry.URL = updateURL
		}
		if cmd.Flags().Changed("notes") {
			entry.Notes = updateNotes
		}
		// (optional: prompt for a new password only if a --password flag was set)

		return svc.UpdateEntry(entry)
	},
}

func init() {
	updateCmd.Flags().StringVar(&updateUsername, "username", "", "new username")
	updateCmd.Flags().StringVar(&updateURL, "url", "", "new URL")
	updateCmd.Flags().StringVar(&updateNotes, "notes", "", "new notes")
	rootCmd.AddCommand(updateCmd)
}