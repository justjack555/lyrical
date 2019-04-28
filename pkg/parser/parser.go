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
	Process a regular file (a song)
 */
func (m *Master) processFile(f *os.File) error {
//	log.Println("ProcessFile(): Reached regular file: ", fInfo.Name())
	sw := newSongWorker(f.Name())
	err := sw.processSong(f)
	if err != nil {
		return err
	}

	// DEBUG
//	log.Println("Master.ProcessFile(): Song words for song: ", f.Name(),
//		":")
	for k, v := range sw.songWords {
//		log.Println("(Word, indices): (", k, ", ", v, ")")
		m.invIndex[k] = append(m.invIndex[k], v)
	}

	return nil
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
func (m *Master) processFileInfo(fInfo os.FileInfo) error {
	f, err := os.Open(fInfo.Name())
	if err != nil {
		log.Println("Master.ProcessFileInfo(): Unable to open file: ", fInfo.Name())
		return err
	}

	defer closeFile(f)

	if !fInfo.IsDir() {
		return m.processFile(f)
	}

//	log.Println("ProcessFileInfo(): File: ", fInfo.Name(), " is a directory...")

	children, err := f.Readdir(0)
	if err != nil {
		return err
	}

//	log.Println("ProcessFileInfo(): Number of children for parent: ", f.Name(), " is ", len(children))

	for _, child := range children {
//		log.Println("ProcessFileInfo(): Processing child file: ", child.Name(), " in parent directory: ", fInfo.Name())
		err = f.Chdir()
		if err != nil {
			return err
		}

		err = m.processFileInfo(child)
		if err != nil {
			log.Println("ProcessFileInfo(): Unable to process file with name: ",
				child.Name(), " in parent dir: ", f.Name())
			return err
		}
	}

	return nil
}

/**
	Driver method to kick off file processing
	Stats the lyrics directory, assures that
	there is no error in doing so, and initiates
	the processing of each file in the directory
	tree
 */
func (m *Master) ProcessFiles(workerChan chan *os.File, done chan bool) error {
	fInfo, err := os.Stat(m.Lyrics_directory)
	if err != nil {
		return err
	}

	err = m.processFileInfo(fInfo)
	if err != nil {
		return err
	}

	log.Println("ProcessFiles(): About to close worker chan: ")
	close(workerChan)
	log.Println("ProcessFiles(): Closed worker chan: ")
	for i := 0; i < m.File_workers; i++ {
		log.Println("ProcessFiles(): Waiting for ", i, "th done ack")
		<- done
	}

	return nil
}

func Start() error {

	m, err := configParser()
	if err != nil {
		return err
	}

	workerChan, done := m.startWorkerPool(m.processFile)

	err = m.ProcessFiles(workerChan, done)
	if err != nil {
		return err
	}

	return nil
}