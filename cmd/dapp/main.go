package main

import (
	"github.com/BigPapaChas/dapp/internal/commands"
	"log"
)

func main() {
	if err := commands.Execute(); err != nil {
		log.Println(err.Error())
	}
}
