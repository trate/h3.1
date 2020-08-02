package main

import (
	"fmt"
	"github.com/trate/h3.1/pkg/transaction"
	"io"
	"log"
	"os"
)

func execute(filename string) (err error) {
	file, err := os.Open(filename)
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

	err = svc.ImportCSV(file)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(svc.Transactions[0])
	return nil
}

func main() {
	if err := execute("export.csv"); err != nil {
		os.Exit(1)
	}
}