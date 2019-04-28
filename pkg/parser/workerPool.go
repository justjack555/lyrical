package parser

import (
	"fmt"
	"log"
	"os"
)

type fileHandler func (file *os.File) (*songWorker, error)

type WorkerChannels struct {
	reqChan chan *os.File
	respChan chan *songWorker
	doneChan chan bool
}

/**
	Each worker waits receives new songs to process from the reqChan
	Once close is sent on the reqChan, the worker acknowledges the
	close on the reqChan by sending true to the doneChan
 */
func startWorker(chans *WorkerChannels, handler fileHandler, i int) {
	logTag := fmt.Sprintf("startWorker%d():", i)

	for v := range chans.reqChan {
		sw, err := handler(v)
		if err != nil {
			log.Println(logTag, "ERR: ", err)
			continue
		}

		chans.respChan <- sw
	}

	chans.doneChan <- true
}

/**
	Intialize each of the channels needed for the worker pool
 */
func createWorkerChans() * WorkerChannels {
	return &WorkerChannels{
		reqChan: make(chan *os.File),
		respChan : make(chan *songWorker),
		doneChan : make(chan bool),
	}
}

/**
	Initialize each of the channels needed for the worker pool
	and prepare each worker to wait for songs to process
 */
func (m *Master) startWorkerPool(handler fileHandler) *WorkerChannels {
	chans := createWorkerChans()

	for i := 0; i < m.File_workers; i++ {
		go startWorker(chans, handler, i)
	}

	return chans
}
