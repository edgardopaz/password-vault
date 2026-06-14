package cli

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"password_vault/internal/vault"
)

var (
	addUsername string
	addURL      string
	addNotes    string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new entry to the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := openVault()
		if err != nil {
			return err
		}

		// the entry's password — read without echo (separate from the master password)
		fmt.Print("Password to store: ")
		pwBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return err
		}

		entry := vault.Entry{
			Username: addUsername,
			Password: string(pwBytes),
			URL:      addURL,
			Notes:    addNotes,
		}
		return svc.AddEntry(entry)
	},
}

func init() {
	addCmd.Flags().StringVar(&addUsername, "username", "", "username for the entry")
	addCmd.Flags().StringVar(&addURL, "url", "", "URL for the entry")
	addCmd.Flags().StringVar(&addNotes, "notes", "", "notes for the entry")
	rootCmd.AddCommand(addCmd)
}