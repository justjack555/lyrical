package main

import (
	"log"
	"github.com/justjack555/lyrical/pkg/parser"
)

func main(){
	log.Println("Initializing parser...")

	err := parser.Start()
	if err != nil {
		log.Fatalln("Parser terminating with error: ", err)
	}
}