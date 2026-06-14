package cli

import (
	"fmt"
	"strconv"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a single entry by id",
	Args:  cobra.ExactArgs(1), // require exactly one positional arg
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid id %q: %w", args[0], err)
		}

		svc, err := openVault()
		if err != nil {
			return err
		}

		entry, err := svc.GetEntry(id)
		if err != nil {
			return err
		}

		// this is the ONE command that reveals the password, by design
		fmt.Printf("Username: %s\nPassword: %s\nURL: %s\nNotes: %s\n",
			entry.Username, entry.Password, entry.URL, entry.Notes)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}