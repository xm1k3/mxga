package utils

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/shopspring/decimal"
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

func ConvertStringToDecimal(str string, decimals int) string {
	d, err := decimal.NewFromString(str)
	if err != nil {
		return "0"
	}

	decimalFactor := decimal.New(1, int32(decimals))

	decimalVal := d.Div(decimalFactor)

	return decimalVal.String()
}
