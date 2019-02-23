package parser

import (
	"github.com/justjack555/lyrical/pkg/common"
	"log"
	"os"
)

/*
	Representation of a word
	in a song.

	We note the song and the position
	of the word in the song
 */
type SongWord struct {
	song *common.Song
	index int
}

/**
	Master coordinates the parsing
	of a set of song into an inverted
	index
 */
type Master struct {
	Song_workers int
	File_workers int
	Lyrics_directory string
	invIndex map[string][]*SongWord
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
	m.invIndex = make(map[string][]*SongWord)

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

func Start() error {

	m, err := configParser()
	if err != nil {
		return err
	}

	log.Println("ConfigParser(): The configured master is: ", *m)

	return nil
}