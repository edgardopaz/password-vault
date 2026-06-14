package cli 

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"password_vault/internal/vault"
	"password_vault/internal/storage"
	"golang.org/x/term"
	"path/filepath"
)

var rootCmd = &cobra.Command {
	Use: "pass",
	Short: "password vault is a cli tool for managing passwords.",
	Long: "password vault is a cli tool for managing passwords - generating, storing, and retrieving them securely",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops, An error while executing Password Vault '%s'\n", err)
		os.Exit(1)
	}
}

func openVault() (*vault.Service, error) {
	// 1. Resolve where the vault file lives (and make sure the folder exists).
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	vaultDir := filepath.Join(configDir, "passvault")
	if err := os.MkdirAll(vaultDir, 0700); err != nil {
		return nil, err
	}
	dbPath := filepath.Join(vaultDir, "vault.db")

	// 2. Open the store at that PATH (this also runs your migrations).
	store, err := storage.NewSQLiteStore(dbPath)
	if err != nil {
		return nil, err
	}

	// 3. Read the master password WITHOUT echoing it to the terminal.
	fmt.Print("Master password: ")
	pwBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // ReadPassword swallows the newline; print one so the next output isn't glued on
	if err != nil {
		return nil, err
	}
	password := string(pwBytes)

	// 4. Build the service. NewService derives the key and, on a fresh DB,
	//    creates+stores the salt (the load-or-create flow you wrote).
	return vault.NewService(store, password)
}