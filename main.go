package main

import (
	"log"

	"github.com/JorritSalverda/infinity/cmd"
)

func main() {
	log.SetFlags(0)
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
