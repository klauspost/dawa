package main

import (
	"fmt"
	"github.com/klauspost/dawa"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	file, err := os.Open("adresser.json") // For read access.
	if err != nil {
		log.Fatal(err)
	}

	iter, err := dawa.ImportAdresserJSON(file)
	if err != nil {
		log.Fatal(err)
	}

	n := 0
	t := time.Now()
	t2 := time.Now()
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
		if n%1000 == 0 {
			fmt.Printf("Processed %d, last 1000 was %.1f per sec.\n", n, (float64)(time.Second*1000)/float64(time.Now().Sub(t2)))
			t2 = time.Now()
		}
		_ = a
		//fmt.Printf("Item:%#v\n", a)
	}
}
