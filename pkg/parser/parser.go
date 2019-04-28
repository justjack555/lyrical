package parser

import (
	"github.com/justjack555/lyrical/pkg/common"
	"log"
	"os"
)

/**
	Master coordinates the parsing
	of a set of song into an inverted
	index
 */
type Master struct {
	Song_workers int
	File_workers int
	Lyrics_directory string
	invIndex map[string][]*songWord
}

type ArgsError struct {
	err string
}

func (ae *ArgsError) Error() string {
	return ae.err
}

/**
	Initialize the new Parsing Master
	by initializing its inverted index
 */
func newMaster() *Master{
	m := new(Master)
	m.invIndex = make(map[string][]*songWord)

	return m
}

/**
	Configure the parser master with the properties
	defined in the configuration file
 */
func configParser() (*Master, error) {
	if numArgs := len(os.Args); numArgs < 2 {
		return nil, &ArgsError{
			err: "Usage: ./parser <config_file>",
		}
	}

	m := newMaster()

	common.LoadConfig(os.Args[1], m)

	return m, nil
}

/**
	Close song lyrics file and check for error
 */
func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal("ProcessFileInfo(): Unable to close file: ", f.Name())
	}
}

/**
	Helper method that determines if a file is a directory
	or a regular file.

	If a regular file, it calls the regular
	file specific processor to parse the song.

	If the file is a directory, the method invokes itself
	on each of the child files

	All errors are returned to the caller for handling
 */
func (m *Master) processFileInfo(fInfo os.FileInfo, reqChan chan *os.File) error {
	f, err := os.Open(fInfo.Name())
	if err != nil {
		log.Println("Master.ProcessFileInfo(): Unable to open file: ", fInfo.Name())
		return err
	}

	if !fInfo.IsDir() {
		reqChan <- f
		return nil
	}

	defer closeFile(f)

	children, err := f.Readdir(0)
	if err != nil {
		return err
	}

	for _, child := range children {
		err = f.Chdir()
		if err != nil {
			return err
		}

		err = m.processFileInfo(child, reqChan)
		if err != nil {
			log.Println("ProcessFileInfo(): Unable to process file with name: ",
				child.Name(), " in parent dir: ", f.Name())
			return err
		}
	}

	return nil
}

/**
	Add the results from a new song to the master
	inverted index
 */
func (m *Master) processResponse(sw *songWorker) {
	for k, v := range sw.songWords {
		m.invIndex[k] = append(m.invIndex[k], v)
	}
}

/**
	Process each of the workers' responses and send acknowledgment
	once all responses have been processed
 */
func (m *Master) processResponses(chans *WorkerChannels) {
	for sw := range chans.respChan {
		m.processResponse(sw)
	}

	chans.doneChan <- true
}

/**
	Driver method to kick off file processing
	Stats the lyrics directory, assures that
	there is no error in doing so, and initiates
	the processing of each file in the directory
	tree
 */
func (m *Master) ProcessFiles(chans *WorkerChannels) error {
	fInfo, err := os.Stat(m.Lyrics_directory)
	if err != nil {
		return err
	}

	go m.processResponses(chans)


	err = m.processFileInfo(fInfo, chans.reqChan)
	if err != nil {
		return err
	}

	close(chans.reqChan)

	for i := 0; i < m.File_workers; i++ {
		<- chans.doneChan
	}

	close(chans.respChan)
	<- chans.doneChan

	return nil
}

/**
	Initialize the pool of worker goroutines
	and launch the processFiles process using
	the worker pool
 */
func Start() error {

	m, err := configParser()
	if err != nil {
		return err
	}

	chans := m.startWorkerPool(processFile)

	err = m.ProcessFiles(chans)
	if err != nil {
		return err
	}

	log.Println("Start(): INFO: Number of words inserted into inverted index: ", len(m.invIndex))

	return nil
}