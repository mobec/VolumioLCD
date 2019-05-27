package utils

import (
	"os"
)

// Exists reports whether the named file or directory exists.
func PathExists(name string) bool {
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}