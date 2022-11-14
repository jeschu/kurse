package support

import (
	"io"
	"log"
)

var Close = func(readCloser io.ReadCloser, errorMsg string) {
	if err := readCloser.Close(); err != nil {
		log.Printf("%s: %v", errorMsg, err)
	}
}
