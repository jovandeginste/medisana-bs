package main

import (
	"encoding/binary"
	"fmt"
	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/examples/lib/dev"
	"golang.org/x/net/context"
	"log"
	"strings"
	"time"
)

func StartBluetooth() {
	d, err := dev.NewDevice(config.Device)
	if err != nil {
		log.Printf("[BLUETOOTH] Can't use device: %s", err)
	}

	ble.SetDefaultDevice(d)
	filter := func(a ble.Advertisement) bool {
		return strings.ToUpper(a.Address().String()) == strings.ToUpper(config.DeviceID)
	}

	time.Sleep(1 * time.Second)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), config.ScanDuration.AsTimeDuration()))
	log.Println("[BLUETOOTH] Starting scan...")
	cln, err := ble.Connect(ctx, filter)

	if err != nil {
		log.Printf("[BLUETOOTH] Scan timeout: %s\n", err)
	} else {
		// Normally, the connection is disconnected by us after our exploration.
		// However, it can be asynchronously disconnected by the remote peripheral.
		// So we wait(detect) the disconnection in the go routine.
		go func() {
			select {
			case <-cln.Disconnected():
				log.Printf("[BLUETOOTH] [ %s ] is disconnected \n", cln.Address())
			case <-time.After(config.Sub.AsTimeDuration()):
				fmt.Println("[BLUETOOTH] [ %s ] timed out\n", cln.Address())
			}
		}()

		log.Printf("[BLUETOOTH] [ %s ] is connected ... Name: '%s'\n", cln.Address(), cln.Name())
		log.Println("[BLUETOOTH] Discovering profile...")
		p, err := cln.DiscoverProfile(true)
		if err != nil {
			log.Printf("[BLUETOOTH] can't discover profile: %s", err)
		} else {

			// Start the exploration.
			explore(cln, p)

			log.Printf("[BLUETOOTH] Discovery done, waiting %d seconds before disconnecting.\n", config.Sub.AsTimeDuration())
			time.Sleep(config.Sub.AsTimeDuration())

			// Disconnect the connection. (On OS X, this might take a while.)
			log.Printf("[BLUETOOTH] Disconnecting [ %s ]... (this might take up to few seconds on OS X)\n", cln.Address())
			cln.ClearSubscriptions()
			cln.CancelConnection()
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
					log.Printf("[BLUETOOTH] subscribe failed: %s\n", err)
				}
			}
		}
	}

	// Then we send the current time (which triggers data transmission)
	for _, s := range p.Services {
		for _, c := range s.Characteristics {
			switch fmt.Sprintf("%s", c.UUID) {
			case "8a81":
				log.Printf("[BLUETOOTH] Sending the time... ")
				thetime := time.Now().Unix()
				binarytime := generateTime(thetime)
				err := cln.WriteCharacteristic(c, binarytime, true)
				if err != nil {
					log.Printf("[BLUETOOTH] Error while writing command: %+v\n", err)
				}
				log.Printf("[BLUETOOTH] done.\n")
			}
		}
	}
	return nil
}

func generateTime(therealtime int64) []byte {
	bs := make([]byte, 4)
	thetime := uint32(therealtime) - 1262304000
	binary.LittleEndian.PutUint32(bs, thetime)
	bs = append([]byte{2}, bs...)

	return bs
}
