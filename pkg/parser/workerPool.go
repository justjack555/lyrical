package parser

import (
	"fmt"
	"log"
	"os"
)

type fileHandler func (file *os.File) error

func startWorker(w chan *os.File, done chan bool, handler fileHandler, i int) {
	logTag := fmt.Sprintf("startWorker%d():", i)
	log.Println(logTag)

	for v := range w {
		log.Println(logTag, " Received val: ", v)
		if err := handler(v); err != nil {
			log.Println(logTag, "ERR: ", err)
		}
	}

	done <- true
}

func (m *Master) startWorkerPool(handler fileHandler) (chan *os.File, chan bool) {
	w := make(chan *os.File)
	done := make(chan bool)

	for i := 0; i < m.File_workers; i++ {
		go startWorker(w, done, handler, i)
	}

	return w, done
}
