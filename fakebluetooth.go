package main

import (
	"bufio"
	"encoding/hex"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

/*
FakeBluetooth fakes receiving data, triggering the plugins

It will pretend to receive data from a scale, which
should trigger all the (configured) plugins.

You can use this to test configurations and/or new plugins
*/
func FakeBluetooth() {
	log.Infoln("[FAKEBLUETOOTH] Sending fake data from 'testdata' to the indicator parser... (waiting 5 seconds)")

	f, err := os.Open("testdata")
	if err != nil {
		log.Fatalf("[FAKEBLUETOOTH] error opening file: %s", err)
	}

	r := bufio.NewReader(f)
	s, e := readln(r)

	time.Sleep(5 * time.Second)

	for e == nil {
		log.Infoln("[FAKEBLUETOOTH] Sending data: ", s)

		h, err := hex.DecodeString(s)
		if err != nil {
			log.Fatalf("[FAKEBLUETOOTH] error decoding line: %v", err)
		}

		go decodeData(h)

		time.Sleep(100 * time.Millisecond)

		s, e = readln(r)
	}

	time.Sleep(1 * time.Second)

	log.Infoln("[FAKEBLUETOOTH] Finished sending fake data from 'testdata' to the indicator parser. Waiting in an infinite loop now.")

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
