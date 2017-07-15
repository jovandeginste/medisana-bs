package main

import (
	"bufio"
	"encoding/hex"
	"log"
	"os"
	"time"
)

/*
FakeBluetooth fakes receiving data, triggering the plugins

It will pretend to receive data from a scale, which
should trigger all the (configured) plugins.

You can use this to test configurations and/or new plugins
*/
func FakeBluetooth() {
	log.Println("[FAKEBLUETOOTH] Sending fake data from 'testdata' to the indicator parser... (waiting 5 seconds)")
	f, err := os.Open("testdata")
	if err != nil {
		log.Println("[FAKEBLUETOOTH] error opening file= ", err)
		os.Exit(1)
	}
	r := bufio.NewReader(f)
	s, e := readln(r)
	time.Sleep(5 * time.Second)
	for e == nil {
		log.Println("[FAKEBLUETOOTH] Sending data: ", s)
		h, _ := hex.DecodeString(s)
		if err != nil {
			log.Printf("[FAKEBLUETOOTH] error decoding line: %+v\n", err)
			os.Exit(1)
		}
		go decodeData(h)
		time.Sleep(200 * time.Millisecond)
		s, e = readln(r)
	}
	time.Sleep(1 * time.Second)
	log.Println("[FAKEBLUETOOTH] Finished sending fake data from 'testdata' to the indicator parser. Waiting in an infinite loop now.")
	select {}
}

func readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
