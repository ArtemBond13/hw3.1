package main

import (
	"github.com/ArtemBond13/hw3.1/pkg/transaction"
	"io"
	"log"
	"os"
)

func main() {
	if err := execute("import.scv"); err != nil {
		os.Exit(1)
	}
}

func execute(filename string) error {
	var err error
	var file *os.File
	if file, err = os.Open(filename); err != nil {
		log.Println(err)
		return err
	}

	defer func(c io.Closer) {
		if cerr := c.Close(); cerr != nil {
			log.Println(cerr)
			if err == nil { // поскольку возвращаемое значение именовано (err),мы в defer имеем к нему доступи
				err = cerr // можем перезаписать значение в нём,но делаем это только если там ещё не было ошибки
			}
		}
	}(file)

	scv := transaction.NewService()

	if err = scv.Import(file); err != nil {
		log.Println(err)
	}
	return err

}
