package main

import (
	"fmt"
	"os"

	"plex-indexer-written-in-go/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
