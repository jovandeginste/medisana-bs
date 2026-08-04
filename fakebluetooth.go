package main

import (
	"bufio"
	"encoding/hex"
	"log/slog"
	"os"
	"time"
)

func fakeBluetoothLogger() *slog.Logger {
	return slog.With("component", "fake-bluetooth")
}

/*
FakeBluetooth fakes receiving data, triggering the plugins

It will pretend to receive data from a scale, which
should trigger all the (configured) plugins.

You can use this to test configurations and/or new plugins
*/
func FakeBluetooth() {
	fakeBluetoothLogger().Info("Sending fake data from 'testdata' to the indicator parser... (waiting 5 seconds)")

	f, err := os.Open("testdata")
	if err != nil {
		fakeBluetoothLogger().Error("error opening file", "err", err)
		os.Exit(1)
	}

	r := bufio.NewReader(f)
	s, e := readln(r)

	time.Sleep(5 * time.Second)

	for e == nil {
		fakeBluetoothLogger().Info("Sending data", "data", s)

		h, err := hex.DecodeString(s)
		if err != nil {
			fakeBluetoothLogger().Error("error decoding line", "err", err)
			os.Exit(1)
		}

		go decodeData(h)

		time.Sleep(100 * time.Millisecond)

		s, e = readln(r)
	}

	time.Sleep(1 * time.Second)

	fakeBluetoothLogger().Info("Finished sending fake data from 'testdata' to the indicator parser. Waiting in an infinite loop now.")

	go debounce()

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
