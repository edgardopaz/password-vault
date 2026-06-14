package cli

import (
	"fmt"
	"strconv"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete an entry by id",
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
		if err := svc.DeleteEntry(id); err != nil {
			return err
		}
		fmt.Printf("Deleted entry %d\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}