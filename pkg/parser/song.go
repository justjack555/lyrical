package parser

import (
	"bufio"
	"github.com/justjack555/lyrical/pkg/common/song"
	"os"
)

/*
	Representation of a word
	in a song.

	We note the song and the sorted position(s)
	of the word in the song
 */
type songWord struct {
	song *song.Song
	indices []int
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

func (sw *songWorker) mapWord(lyric string, i int) {
	res := sw.songWords[lyric]
	if res == nil {
		res = new(songWord)
		res.song = sw.song
	}

	res.indices = append(res.indices, i)
	sw.songWords[lyric] = res
}

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

func (sw *songWorker) processSong(f *os.File) error {

	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanWords)

	err := sw.mapWords(sc)
	if err != nil {
		return err
	}

	return nil
}