package main

import (
	"log"
	"github.com/justjack555/lyrical/pkg/parser"
)

func main(){
	log.Println("Initializing parser...")

	m, err := parser.Start()
	if err != nil {
		log.Fatalln("Parser terminating with error: ", err)
		return
	}

	err = m.StoreIndex()
	if err != nil {
		log.Fatalln("Parser terminating with error: ", err)
		return
	}
}