package parser

import (
	"bufio"
	"github.com/justjack555/lyrical/pkg/common/song"
	"os"
)

/*
	Representation of a word
	in a Song.

	We note the Song and the sorted position(s)
	of the word in the Song
 */
type songWord struct {
	Song    *song.Song
	Indices []int
}

type songWorker struct {
	songWords map[string]*songWord
	song *song.Song
}

func newSongWorker(name string) *songWorker {
	sw := new(songWorker)
	sw.songWords = make(map[string]*songWord)
	sw.song = song.New(name)
	return sw
}

/**
	Either create or update the representation of a
	word's appearance in a particular Song
 */
func (sw *songWorker) mapWord(lyric string, i int) {
	res := sw.songWords[lyric]
	if res == nil {
		res = new(songWord)
		res.Song = sw.song
	}

	res.Indices = append(res.Indices, i)
	sw.songWords[lyric] = res
}

/**
	Obtain the index position of a word and its text, and add that to
	the representation of the document
 */
func (sw *songWorker) mapWords(sc *bufio.Scanner) error {
	for i := 0; sc.Scan(); i++ {
		lyric := sc.Text()
		sw.mapWord(lyric, i)
	}

	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}

/**
	Setup buffered scanner to parse words from Song's lyric file,
	then invoke the word mapper
 */
func (sw *songWorker) processSong(f *os.File) error {

	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanWords)

	err := sw.mapWords(sc)
	if err != nil {
		return err
	}

	return nil
}

/**
	Process a regular file (a Song) by initializing a new Song worker,
	having it process the Song's lyric file, and returning th eresult
 */
func processFile(f *os.File) (*songWorker, error) {
	//	log.Println("ProcessFile(): Reached regular file: ", fInfo.Name())
	defer closeFile(f)

	sw := newSongWorker(f.Name())
	err := sw.processSong(f)
	if err != nil {
		return nil, err
	}

	return sw, nil
}