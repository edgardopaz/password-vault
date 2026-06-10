package cli 

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command {
	Use: "pass",
	Short: "password vault is a cli tool for managing passwords.",
	Long: "password vault is a cli tool for managing passwords - generating, storing, and retrieving them securely",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops, An error while executing Password Vault '%s'\n", err)
		os.Exit(1)
	}
}

func openVault() (*vault.Service, error) {
	// resolve db path
	
	store, err := storage.NewSQLiteStore(password)
	if err != nil {
		return err
	}
	return nil
	password := os.Getenv("VAULT_PASSWORD")
	if password == "" {
		return nil, errors.New("vault password is not set")
	}
	return vault.NewService(store, password)
}