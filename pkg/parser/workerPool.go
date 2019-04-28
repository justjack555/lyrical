package parser

import (
	"fmt"
	"log"
	"os"
)

type fileHandler func (file *os.File) error

type WorkerChannels struct {
	reqChan chan *os.File
	respChan chan *songWorker
	doneChan chan bool
}

func startWorker(chans *WorkerChannels, handler fileHandler, i int) {
	logTag := fmt.Sprintf("startWorker%d():", i)
	log.Println(logTag)

	for v := range chans.reqChan {
		log.Println(logTag, " Received val: ", v)
		if err := handler(v); err != nil {
			log.Println(logTag, "ERR: ", err)
		}
	}

	chans.doneChan <- true
}

func createWorkerChans() * WorkerChannels {
	return &WorkerChannels{
		reqChan: make(chan *os.File),
		respChan : make(chan *songWorker),
		doneChan : make(chan bool),
	}
}

func (m *Master) startWorkerPool(handler fileHandler) *WorkerChannels {
	chans := createWorkerChans()

	for i := 0; i < m.File_workers; i++ {
		go startWorker(chans, handler, i)
	}

	return chans
}
