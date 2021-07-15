package utils

import (
	"encoding/csv"
	"log"
	"os"
)

func WriteCSV(name string, record []string) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	writer := csv.NewWriter(file)
	err = writer.Write(record)
	if err != nil {
		log.Fatal(err)
	}
	writer.Flush()
}
