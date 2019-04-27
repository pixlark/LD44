package main

import (
	"fmt"
	"os"
)

func fatal(message string) {
	fmt.Fprintf(os.Stderr, "Fatal error: %s", message)
	os.Exit(1)
}
