package main

import (
	"fmt"
	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/examples/lib/dev"
	"golang.org/x/net/context"
	"log"
	"strings"
	"time"
)

func StartBluetooth() {
	d, err := dev.NewDevice(device)
	if err != nil {
		log.Printf("Can't use new device: %s", err)
	}
	ble.SetDefaultDevice(d)

	filter := func(a ble.Advertisement) bool {
		return strings.ToUpper(a.Address().String()) == strings.ToUpper(deviceID)
	}

	for {
		ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), scanDuration))
		cln, err := ble.Connect(ctx, filter)
		if err != nil {
			log.Printf("Timeout: %s\n", err)
		} else {
			// Make sure we had the chance to print out the message.
			done := make(chan struct{})
			// Normally, the connection is disconnected by us after our exploration.
			// However, it can be asynchronously disconnected by the remote peripheral.
			// So we wait(detect) the disconnection in the go routine.
			go func() {
				<-cln.Disconnected()
				log.Printf("[ %s ] is disconnected \n", cln.Address())
				close(done)
			}()

			log.Printf("Discovering profile...\n")
			p, err := cln.DiscoverProfile(true)
			if err != nil {
				log.Printf("can't discover profile: %s", err)
			}

			// Start the exploration.
			explore(cln, p)

			time.Sleep(sub)

			// Disconnect the connection. (On OS X, this might take a while.)
			log.Printf("Disconnecting [ %s ]... (this might take up to few seconds on OS X)\n", cln.Address())
			cln.CancelConnection()

			<-done
		}
	}
}

func explore(cln ble.Client, p *ble.Profile) error {
	// First we subscribe to the metrics

	for _, s := range p.Services {
		for _, c := range s.Characteristics {
			switch fmt.Sprintf("%s", c.UUID) {
			case "8a21", "8a22", "8a82":
				h := func(req []byte) { decodeData(req) }

				if err := cln.Subscribe(c, true, h); err != nil {
					log.Printf("subscribe failed: %s\n", err)
				}
			}
		}
	}

	// Then we send the current time (which triggers data transmission)
	for _, s := range p.Services {
		for _, c := range s.Characteristics {
			switch fmt.Sprintf("%s", c.UUID) {
			case "8a81":
				log.Printf("Sending the time... ")
				thetime := time.Now().Unix()
				binarytime := generateTime(thetime)
				err := cln.WriteCharacteristic(c, binarytime, true)
				if err != nil {
					log.Printf("Error while writing command: %+v\n", err)
				}
				log.Printf("done.\n")
			}
		}
	}
	return nil
}
