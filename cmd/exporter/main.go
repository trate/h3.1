package main

import (
	"github.com/trate/h3.1/pkg/transaction"
	"io"
	"log"
	"os"
)

func execute(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(c io.Closer) {
		if cerr := c.Close(); cerr != nil {
			log.Println(cerr)
			if err == nil {
				err = cerr
			}
		}
	}(file)

	svc := transaction.NewService()

	_, err = svc.Register("0001", "0002", 10_000_000)
	if err != nil {
		log.Println(err)
		return
	}

	err = svc.ExportJSON(file)
	if err != nil {
		log.Println(err)
		return
	}
	return nil
}

func main() {
	if err := execute("export.json"); err != nil {
		os.Exit(1)
	}
}