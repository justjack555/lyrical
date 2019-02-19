package main

import (
	"log"
	"lyrical/pkg/parser"
)

func main(){
	log.Println("Initializing parser...")

	err := parser.Start()
	if err != nil {
		log.Fatalln("Parser terminating with error: ", err)
	}
}