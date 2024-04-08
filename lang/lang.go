package lang

import (
	"io"
	"log"
)

var Close = func(readCloser io.ReadCloser, errorMsg string) {
	if err := readCloser.Close(); err != nil {
		log.Printf("%s: %v", errorMsg, err)
	}
}

func FatalOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
