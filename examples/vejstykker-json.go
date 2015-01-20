package main

import (
	"github.com/klauspost/dawa"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	file, err := os.Open("vejstykker.json") // For read access.
	if err != nil {
		log.Fatal(err)
	}

	iter, err := aws.ImportVejstykkerJSON(file)
	if err != nil {
		log.Fatal(err)
	}

	n := 0
	t := time.Now()
	for {
		a, err := iter.Next()
		if err == io.EOF {
			log.Printf("Finished reading %d entries in %v.\n", n, time.Now().Sub(t))
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		n++
		log.Printf("Entry:%#v\n", a)
	}
}
