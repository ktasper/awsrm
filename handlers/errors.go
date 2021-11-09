package handlers

import (
	"fmt"
	"os"
)

// Error handling
func ExitErrorf(msg string, args ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, msg+"\n", args...)
	if err != nil {
		return
	}
	os.Exit(1)
}
