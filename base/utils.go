package base

import (
	"github.com/jouir/pgterminate/log"
)

// Panic prints a non-nil error and terminates the program
func Panic(err error) {
	if err != nil {
		log.Fatalf("%s\n", err)
	}
}

// InSlice detects value presence in a string slice
func InSlice(value string, slice []string) bool {
	for _, val := range slice {
		if value == val {
			return true
		}
	}
	return false
}
