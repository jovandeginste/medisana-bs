package main

import (
	"fmt"
	"github.com/currantlabs/ble"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

func explore(cln ble.Client, p *ble.Profile) error {
	// First we subscribe to the metrics

	for _, s := range p.Services {
		for _, c := range s.Characteristics {
			fmt.Printf("Indication characteristic: %s\n", c)
			switch fmt.Sprintf("%s", c.UUID) {
			case "8a21", "8a22", "8a82":
				fmt.Printf("\n-- Subscribe to indication of %v --\n", c.UUID)

				h := func(req []byte) { parseIndication(req) }

				if err := cln.Subscribe(c, true, h); err != nil {
					fmt.Printf("subscribe failed: %s\n", err)
				}
			}
		}
	}

	// Then we send the current time (which triggers data transmission)
	for _, s := range p.Services {
		for _, c := range s.Characteristics {
			fmt.Printf("Indication characteristic: %s\n", c)
			switch fmt.Sprintf("%s", c.UUID) {
			case "8a81":
				for _, d := range c.Descriptors {
					fmt.Printf("Descriptor: %+v\n", d)
				}

				fmt.Printf("Sending the time... ")
				thetime := time.Now().Unix()
				binarytime := generateTime(thetime)
				err := cln.WriteCharacteristic(c, binarytime, true)
				if err != nil {
					fmt.Printf("Error while writing command: %+v\n", err)
				}
				fmt.Printf("done.\n")
			}
		}
	}
	return nil
}

func propString(p ble.Property) string {
	var s string
	for k, v := range map[ble.Property]string{
		ble.CharBroadcast:   "B",
		ble.CharRead:        "R",
		ble.CharWriteNR:     "w",
		ble.CharWrite:       "W",
		ble.CharNotify:      "N",
		ble.CharIndicate:    "I",
		ble.CharSignedWrite: "S",
		ble.CharExtended:    "E",
	} {
		if p&k != 0 {
			s += v
		}
	}
	return s
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		fmt.Printf("canceled\n")
	default:
		fmt.Printf(err.Error())
	}
}
func parseIndication(req []byte) {
	fmt.Printf("Got data: [% X]\n", req)
	switch req[0] {
	case 0x84:
		fmt.Printf("Received person data: ")
		person := decodePerson(req)
		fmt.Printf("%+v\n", person)
	case 0x1D:
		fmt.Printf("Received weight data: ")
		weight := decodeWeight(req)
		fmt.Printf("%+v\n", weight)
	case 0x6F:
		fmt.Printf("Received body data: ")
		body := decodeBody(req)
		fmt.Printf("%+v\n", body)
	default:
		fmt.Println("Unhandled Indication encountered")
	}
}
