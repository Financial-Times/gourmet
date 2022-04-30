package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gotha/gourmet/cmd"
)

func main() {
	r := cmd.NewRouter()

	RegisterTransport(r)

	err := r.Run(context.Background())
	if err != nil {
		if err == cmd.ErrNotFound {
			fmt.Println("Usage: <cmd> get user (type) (flags)")
			fmt.Println("\t data --user-id <num>")
			fmt.Println("\t profile --email <email> --full")
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
