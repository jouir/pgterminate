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
