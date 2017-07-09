package main

import (
	"bufio"
	"encoding/hex"
	"log"
	"os"
	"time"
)

func FakeBluetooth() {
	log.Println("Sending fake data from 'testdata' to the indicator parser... (waiting 5 seconds)")
	f, err := os.Open("testdata")
	if err != nil {
		log.Println("error opening file= ", err)
		os.Exit(1)
	}
	r := bufio.NewReader(f)
	s, e := Readln(r)
	time.Sleep(5 * time.Second)
	for e == nil {
		log.Println("Sending data: ", s)
		h, _ := hex.DecodeString(s)
		if err != nil {
			log.Printf("error decoding line: %+v\n", err)
			os.Exit(1)
		}
		decodeData(h)
		time.Sleep(200 * time.Millisecond)
		s, e = Readln(r)
	}
	time.Sleep(1 * time.Second)
	log.Println("Finished sending fake data from 'testdata' to the indicator parser. Waiting in an infinite loop now.")
	select {}
}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
