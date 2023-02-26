package utils

import (
	"io/ioutil"
	"log"
	"os"
)

func ReadPrivateKey(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
