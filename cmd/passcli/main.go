package main

import (
	"password_vault/internal/cli"
	"flag"
)

func main() {
	add := flag.String("add", "", "Add a new password to the vault")
	get := flag.String("get", "", "Get a password from the vault")
	list := flag.String("list", "", "List all passwords in the vault")
	
	flag.Parse()
	cli.Execute()
}