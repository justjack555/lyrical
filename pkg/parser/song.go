package parser

import (
	"bufio"
	"log"
	"os"
)

type songWorker struct {
	songWords map[string]*SongWord
}

func newSongWorker() *songWorker {
	sw := new(songWorker)
	sw.songWords = make(map[string]*SongWord)
	return sw
}

func (sw *songWorker) mapWords(sc *bufio.Scanner) error {
	for sc.Scan() {
		lyric := sc.Text()
		log.Println("Song.MapWords(): Lyric: ", lyric)
	}

	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}

func (sw *songWorker) processSong(f *os.File) error {

	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanWords)

	err := sw.mapWords(sc)
	if err != nil {
		return err
	}
	return nil
}