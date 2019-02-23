package common

import (
	"bufio"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
)

/**
	Read all bytes from file of length size into
	a byte array
 */
func getBytesFromReader(br *bufio.Reader, size int64) []byte {
	b := make([]byte, size)
	for {
		_, err := br.Read(b)
		if err != nil {
			if err != io.EOF {
				log.Fatalln("ERR: While reading from file: ", err)
			}
			break;
		}
	}

	//	log.Println("LOAD_CONFIG(): Byte slice read from file is: ", string(b))
	return b
}

/**
	Read the server configuration for the
	master from config file
 */
func LoadConfig(path string, c interface{}) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("ERR: Unable to read config file. Terminating...")
	}
	defer f.Close()

	fileinfo, err := f.Stat()
	if err != nil {
		log.Fatalln("ERR: Unable to stat config file info: ", err)
	}

	filesize := fileinfo.Size()
	br := bufio.NewReader(f)
	b := getBytesFromReader(br, filesize)

	err = yaml.Unmarshal(b, c)
	if err != nil {
		log.Fatal("ERR: ", err)
	}
}