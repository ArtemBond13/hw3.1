package main

import (
	"github.com/ArtemBond13/hw3.1/pkg/transaction"
	"io"
	"log"
	"os"
)

func main() {
	if err := execute("export.xml"); err != nil {
		os.Exit(1)
	}

}
func execute(filename string) error {
	var err error
	var file *os.File
	if file, err = os.Create(filename); err != nil {
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

	if _, err := scv.Register("0001", "0002", 100_000_00); err != nil {
		log.Println(err)
	}
	if _, err := scv.Register("0003", "0004", 100_000_00); err != nil {
		log.Println(err)
	}

	if err = scv.ExportXML("export.xml"); err != nil {
		log.Println(err)
	}
	return err

}
