package main

import (
	"fmt"
	"os"
)

func fatal(message string) {
	fmt.Fprintf(os.Stderr, "Fatal error: %s\n", message)
	os.Exit(1)
}
