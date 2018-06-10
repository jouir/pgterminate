package base

import (
	"log"
)

// Panic prints a non-nil error and terminates the program
func Panic(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
